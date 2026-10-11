package server

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"trae2api-web/internal/auth"
	"trae2api-web/internal/pool"
)

func (h *Handler) adminKeyInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"configured": h.cfg.APIKey != ""})
}

func (h *Handler) adminKeyCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) compatHealth(w http.ResponseWriter, r *http.Request) {
	statuses := h.cfg.Pool.List()
	healthy := 0
	for _, s := range statuses {
		if s.Enabled && !s.Disabled && !s.Cooling {
			healthy++
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok", "available_accounts": healthy, "total_accounts": len(statuses),
		"available_providers": boolInt(healthy > 0), "total_providers": 1,
	})
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func (h *Handler) compatAccounts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"accounts": h.compatAccountList()})
}

func (h *Handler) compatAccountList() []map[string]any {
	statuses := h.cfg.Pool.List()
	out := make([]map[string]any, 0, len(statuses))
	for _, s := range statuses {
		out = append(out, h.compatAccount(s))
	}
	return out
}

func (h *Handler) compatAccount(s pool.Status) map[string]any {
	healthy := s.Enabled && !s.Disabled && !s.Cooling
	item := map[string]any{
		"uid": s.UID, "provider": "trae", "nickname": s.Nickname, "enabled": s.Enabled,
		"disabled": s.Disabled, "cooling": s.Cooling, "healthy": healthy,
		"reason": s.Reason, "credits_remaining": s.Credits, "credits_total": nil,
		"credits_used": nil, "credits_packs": 0, "failure_count": s.ErrCount,
		"cooldown_until": nil, "profile": "trae-solo", "site": "TRAE",
	}
	if !s.Until.IsZero() {
		item["cooldown_until"] = s.Until.Unix()
	}
	if a := h.cfg.Pool.AuthByUID(s.UID); a != nil {
		item["expires_at"] = a.ExpiresAt
		item["expired_soon"] = a.ExpiresAt > 0 && a.NeedsRefresh(24*time.Hour)
		item["machine_id"] = prefix(a.MachineID, 8)
		item["device_id"] = prefix(a.DeviceID, 8)
	}
	return item
}

func (h *Handler) compatOverview(w http.ResponseWriter, r *http.Request) {
	accounts := h.compatAccountList()
	models := h.modelList()
	var remaining int64
	for _, a := range accounts {
		if value, ok := a["credits_remaining"].(int64); ok {
			remaining += value
		}
	}
	usage := h.usage.summary()
	writeJSON(w, http.StatusOK, map[string]any{
		"accounts": accounts, "models": models, "model_count": len(models),
		"credits": map[string]any{"remaining": remaining, "total": remaining, "used": int64(0)},
		"usage":   usage,
	})
}

func (h *Handler) compatEnableAccount(w http.ResponseWriter, r *http.Request) {
	h.compatSetEnabled(w, r, true)
}

func (h *Handler) compatDisableAccount(w http.ResponseWriter, r *http.Request) {
	h.compatSetEnabled(w, r, false)
}

func (h *Handler) compatSetEnabled(w http.ResponseWriter, r *http.Request, enabled bool) {
	uid := r.PathValue("uid")
	if !h.cfg.Pool.SetEnabled(uid, enabled, map[bool]string{true: "", false: "user disabled"}[enabled]) {
		writeOpenAIError(w, http.StatusNotFound, "not_found", "account not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "enabled": enabled})
}

func (h *Handler) compatAliasAccount(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	var req struct {
		Alias string `json:"alias"`
	}
	if err := decodeBody(w, r, &req); err != nil {
		return
	}
	a := h.cfg.Pool.AuthByUID(uid)
	if a == nil {
		writeOpenAIError(w, http.StatusNotFound, "not_found", "account not found")
		return
	}
	a.Nickname = strings.TrimSpace(req.Alias)
	if a.FilePath == "" {
		a.FilePath = auth.FilePathFor(h.cfg.AuthDir, uid)
	}
	if err := a.SaveAtomic(); err != nil {
		writeOpenAIError(w, http.StatusInternalServerError, "save_failed", err.Error())
		return
	}
	h.cfg.Pool.Add(a)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "alias": a.Nickname})
}

func (h *Handler) compatDeleteAccount(w http.ResponseWriter, r *http.Request) {
	h.adminDeleteAccount(w, r)
}

type compatCredit struct {
	UID            string `json:"uid"`
	Nickname       string `json:"nickname"`
	Remain         int64  `json:"remain"`
	Limit          int64  `json:"limit"`
	Used           int64  `json:"used"`
	Packs          int    `json:"packs"`
	CheckedIn      bool   `json:"checked_in"`
	CheckinCredits int64  `json:"checkin_credits"`
	CheckinEnable  bool   `json:"checkin_enable"`
	Cooling        bool   `json:"cooling"`
	Disabled       bool   `json:"disabled"`
	Error          string `json:"error,omitempty"`
}

func (h *Handler) collectCompatCredits() []compatCredit {
	statuses := h.cfg.Pool.List()
	out := make([]compatCredit, len(statuses))
	var wg sync.WaitGroup
	for i, status := range statuses {
		wg.Add(1)
		go func(i int, status pool.Status) {
			defer wg.Done()
			item := compatCredit{UID: status.UID, Nickname: status.Nickname, Cooling: status.Cooling, Disabled: status.Disabled, Remain: status.Credits}
			a := h.cfg.Pool.AuthByUID(status.UID)
			if a == nil || h.cfg.Upstream == nil {
				item.Error = "no auth or upstream"
				out[i] = item
				return
			}
			remain, limit, used, packs, err := h.cfg.Upstream.EntUsage(a)
			if err != nil {
				item.Error = "ent_usage: " + err.Error()
			} else {
				item.Remain, item.Limit, item.Used, item.Packs = remain, limit, used, packs
				h.cfg.Pool.SetCredits(status.UID, remain)
			}
			checked, checkCredits, enable, checkErr := h.cfg.Upstream.CheckinStatus(a)
			if checkErr != nil {
				if item.Error != "" {
					item.Error += "; "
				}
				item.Error += "checkin: " + checkErr.Error()
			} else {
				item.CheckedIn, item.CheckinCredits, item.CheckinEnable = checked, checkCredits, enable
			}
			out[i] = item
		}(i, status)
	}
	wg.Wait()
	return out
}

func (h *Handler) compatRefreshCredits(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "accounts": h.collectCompatCredits()})
}

func (h *Handler) compatCheckin(w http.ResponseWriter, r *http.Request) {
	for _, status := range h.cfg.Pool.List() {
		if status.Disabled {
			continue
		}
		a := h.cfg.Pool.AuthByUID(status.UID)
		if a == nil || h.cfg.Upstream == nil {
			continue
		}
		checked, _, enable, err := h.cfg.Upstream.CheckinStatus(a)
		if err == nil && !checked && enable {
			_ = h.cfg.Upstream.CheckinClaim(a)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "accounts": h.collectCompatCredits()})
}

func (h *Handler) compatModels(w http.ResponseWriter, r *http.Request) {
	models, source := h.compatModelList()
	writeJSON(w, http.StatusOK, map[string]any{"models": models, "source": source})
}

func (h *Handler) compatModelCatalog(w http.ResponseWriter, r *http.Request) {
	models, _ := h.compatModelList()
	writeJSON(w, http.StatusOK, map[string]any{"models": models})
}

func (h *Handler) compatModelList() ([]map[string]any, string) {
	models := h.modelList()
	dynamicModelsCache.RLock()
	dynamic := len(dynamicModelsCache.ids) > 0
	dynamicModelsCache.RUnlock()
	if dynamic {
		return models, "dynamic"
	}
	return models, "static"
}

func (h *Handler) compatRefreshModels(w http.ResponseWriter, r *http.Request) {
	dynamicModelsCache.Lock()
	dynamicModelsCache.ids = nil
	dynamicModelsCache.fetched = time.Time{}
	dynamicModelsCache.lastFail = time.Time{}
	dynamicModelsCache.Unlock()
	models, source := h.compatModelList()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "models": models, "source": source})
}

func (h *Handler) compatProviders(w http.ResponseWriter, r *http.Request) {
	health := h.compatHealthSnapshot()
	writeJSON(w, http.StatusOK, map[string]any{"providers": []map[string]any{{
		"name": "trae", "display_name": "TRAE SOLO", "ready": health.healthy > 0, "default": true,
		"capabilities": []string{"models", "credits", "checkin", "oauth", "upload", "add_account"},
		"status":       map[string]any{"healthy_count": health.healthy, "account_count": health.total, "model_count": len(h.modelList())},
	}}})
}

func (h *Handler) compatProvider(w http.ResponseWriter, r *http.Request) {
	if r.PathValue("name") != "trae" {
		writeOpenAIError(w, http.StatusNotFound, "not_found", "provider not found")
		return
	}
	models, source := h.compatModelList()
	health := h.compatHealthSnapshot()
	provider := map[string]any{"name": "trae", "display_name": "TRAE SOLO", "ready": health.healthy > 0, "default": true, "capabilities": []string{"models", "credits", "checkin", "oauth", "upload", "add_account"}, "status": map[string]any{"healthy_count": health.healthy, "account_count": health.total, "model_count": len(models), "model_source": source}}
	writeJSON(w, http.StatusOK, map[string]any{"provider": provider, "accounts": h.compatAccountList(), "models": models})
}

type compatHealthResult struct{ healthy, total int }

func (h *Handler) compatHealthSnapshot() compatHealthResult {
	result := compatHealthResult{}
	for _, s := range h.cfg.Pool.List() {
		result.total++
		if s.Enabled && !s.Disabled && !s.Cooling {
			result.healthy++
		}
	}
	return result
}

func (h *Handler) compatUsageSummary(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.usage.summary())
}

func (h *Handler) compatUsageTimeseries(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.usage.timeseries()})
}

func (h *Handler) compatUsageRecent(w http.ResponseWriter, r *http.Request) {
	records := h.usage.snapshot()
	writeJSON(w, http.StatusOK, map[string]any{"records": records, "total": len(records)})
}

func (h *Handler) compatSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.compatSettingsValue())
}

func (h *Handler) compatSettingsValue() map[string]any {
	return map[string]any{
		"listen": h.cfg.Listen, "callback_port": h.cfg.CallbackPort, "auth_dir": h.cfg.AuthDir,
		"state_file": h.cfg.StateFile, "default_model": h.cfg.DefaultModel,
		"api_key_configured": h.cfg.APIKey != "", "plan_credit": h.cfg.PlanCooldown.String(),
		"soft_rate": h.cfg.SoftCooldown.String(), "err_threshold": h.cfg.ErrThreshold,
		"err_cooldown": h.cfg.ErrCooldown.String(), "upstream_timeout_seconds": h.cfg.UpstreamTimeoutSeconds,
	}
}

func (h *Handler) compatSaveSettings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DefaultModel string `json:"default_model"`
	}
	if err := decodeBodyOptional(r, &req); err != nil {
		writeOpenAIError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if strings.TrimSpace(req.DefaultModel) != "" {
		h.cfg.DefaultModel = strings.TrimSpace(req.DefaultModel)
	}
	data := h.compatSettingsValue()
	data["ok"] = true
	writeJSON(w, http.StatusOK, data)
}
