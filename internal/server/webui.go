package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

// adminUIFS 是 Vue 构建产物。前端源码在 webui/，构建后由 scripts/sync-embed.mjs
// 同步到这里；单个 Go 二进制即可提供完整管理台。
//
//go:embed all:admin-ui
var adminUIFS embed.FS

func (h *Handler) adminUIPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	data, err := fs.ReadFile(adminUIFS, "admin-ui/index.html")
	if err != nil {
		http.Error(w, "admin UI is not built", http.StatusNotFound)
		return
	}
	_, _ = w.Write(data)
}

func (h *Handler) adminUIHandler() http.Handler {
	sub, err := fs.Sub(adminUIFS, "admin-ui")
	if err != nil {
		return http.NotFoundHandler()
	}
	assets := http.StripPrefix("/admin-ui/", http.FileServer(http.FS(sub)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		if strings.HasSuffix(r.URL.Path, ".html") || r.URL.Path == "/admin-ui/" {
			w.Header().Set("Cache-Control", "no-store")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		assets.ServeHTTP(w, r)
	})
}
