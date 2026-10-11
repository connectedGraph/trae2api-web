package server

import (
	"encoding/json"
	"sort"
	"sync"
	"time"
)

// usageTracker 保存当前进程内的轻量调用统计，供管理台展示。
// 凭证、请求正文和模型输出不会写入这里；服务重启后统计从零开始。
type usageTracker struct {
	mu      sync.RWMutex
	nextID  int64
	records []usageRecord
}

type usageRecord struct {
	ID           int64     `json:"id"`
	Timestamp    time.Time `json:"-"`
	TS           int64     `json:"ts"`
	Model        string    `json:"model"`
	Protocol     string    `json:"protocol"`
	AccountUID   string    `json:"account_uid"`
	InputTokens  int64     `json:"input_tokens"`
	OutputTokens int64     `json:"output_tokens"`
	TotalTokens  int64     `json:"total_tokens"`
	Status       string    `json:"status"`
	LatencyMS    int64     `json:"latency_ms"`
}

func newUsageTracker() *usageTracker { return &usageTracker{records: make([]usageRecord, 0, 200)} }

func (u *usageTracker) record(model, account, protocol string, tokens int64, ok bool, started time.Time) {
	if u == nil {
		return
	}
	if tokens < 0 {
		tokens = 0
	}
	status := "error"
	if ok {
		status = "ok"
	}
	now := time.Now()
	u.mu.Lock()
	defer u.mu.Unlock()
	u.nextID++
	u.records = append(u.records, usageRecord{
		ID: u.nextID, Timestamp: now, TS: now.Unix(), Model: model, Protocol: protocol,
		AccountUID: account, TotalTokens: tokens, Status: status,
		LatencyMS: time.Since(started).Milliseconds(),
	})
	if len(u.records) > 200 {
		u.records = u.records[len(u.records)-200:]
	}
}

func (u *usageTracker) snapshot() []usageRecord {
	if u == nil {
		return nil
	}
	u.mu.RLock()
	defer u.mu.RUnlock()
	out := append([]usageRecord(nil), u.records...)
	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp.After(out[j].Timestamp) })
	return out
}

func (u *usageTracker) summary() map[string]any {
	records := u.snapshot()
	startToday := time.Now().Truncate(24 * time.Hour)
	// Truncate 使用 UTC；按本地日期重新计算当天起点，避免跨时区时概览误差。
	now := time.Now()
	startToday = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	byModel := map[string]*usageAggregate{}
	byAccount := map[string]*usageAggregate{}
	var total, today, totalTokens, todayTokens int64
	for _, r := range records {
		total++
		totalTokens += r.TotalTokens
		if !r.Timestamp.Before(startToday) {
			today++
			todayTokens += r.TotalTokens
		}
		m := byModel[r.Model]
		if m == nil {
			m = &usageAggregate{}
			byModel[r.Model] = m
		}
		m.Count++
		m.Tokens += r.TotalTokens
		a := byAccount[r.AccountUID]
		if a == nil {
			a = &usageAggregate{}
			byAccount[r.AccountUID] = a
		}
		a.Count++
		a.Tokens += r.TotalTokens
	}
	models := make([]map[string]any, 0, len(byModel))
	for model, v := range byModel {
		models = append(models, map[string]any{"model": model, "count": v.Count, "tokens": v.Tokens})
	}
	accounts := make([]map[string]any, 0, len(byAccount))
	for uid, v := range byAccount {
		accounts = append(accounts, map[string]any{"account_uid": uid, "count": v.Count, "tokens": v.Tokens})
	}
	sort.Slice(models, func(i, j int) bool { return models[i]["count"].(int64) > models[j]["count"].(int64) })
	sort.Slice(accounts, func(i, j int) bool { return accounts[i]["count"].(int64) > accounts[j]["count"].(int64) })
	return map[string]any{
		"total_requests": total, "total_tokens": totalTokens,
		"today_requests": today, "today_tokens": todayTokens,
		"by_model": models, "by_account": accounts,
	}
}

type usageAggregate struct{ Count, Tokens int64 }

func (u *usageTracker) timeseries() []map[string]any {
	now := time.Now()
	start := now.Truncate(time.Hour).Add(-23 * time.Hour)
	buckets := make([]map[string]any, 24)
	for i := range buckets {
		ts := start.Add(time.Duration(i) * time.Hour)
		buckets[i] = map[string]any{"bucket": ts.Format("01-02 15:00"), "bucket_ts": ts.Unix(), "count": int64(0), "tokens": int64(0)}
	}
	for _, r := range u.snapshot() {
		if r.Timestamp.Before(start) {
			continue
		}
		idx := int(r.Timestamp.Truncate(time.Hour).Sub(start) / time.Hour)
		if idx < 0 || idx >= len(buckets) {
			continue
		}
		buckets[idx]["count"] = buckets[idx]["count"].(int64) + 1
		buckets[idx]["tokens"] = buckets[idx]["tokens"].(int64) + r.TotalTokens
	}
	return buckets
}

func tokensFromResponse(resp map[string]any) int64 {
	usage, _ := resp["usage"].(map[string]any)
	if usage == nil {
		return 0
	}
	if value, ok := usage["total_tokens"]; ok {
		switch v := value.(type) {
		case float64:
			return int64(v)
		case int64:
			return v
		case int:
			return int64(v)
		case json.Number:
			n, _ := v.Int64()
			return n
		}
	}
	return 0
}
