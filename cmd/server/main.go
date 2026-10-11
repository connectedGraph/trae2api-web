// main.go trae2api-web 入口：加载配置 → 构建 pool → 起 HTTP 服务。
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trae2api-web/internal/auth"
	"trae2api-web/internal/pool"
	"trae2api-web/internal/scheduler"
	"trae2api-web/internal/server"
	"trae2api-web/internal/upstream"
)

func main() {
	cfgPath := flag.String("config", "config.json", "path to config json")
	flag.Parse()

	cfg, err := Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	auths, err := auth.LoadDir(cfg.AuthDir)
	if err != nil {
		log.Fatalf("load auths: %v", err)
	}
	log.Printf("loaded %d account(s) from %s", len(auths), cfg.AuthDir)

	p := pool.New(cfg.StateFile)
	p.SyncToDir(auths) // 对齐：剔除 state.json 中已删除 auth 文件的幽灵账号

	up := upstream.New()
	up.HTTP.Timeout = time.Duration(cfg.Upstream.TimeoutSeconds) * time.Second
	// 流式客户端无总超时，仅用首字节兜底（时长由 SSE 流本身决定）。
	if tr, ok := up.StreamHTTP.Transport.(*http.Transport); ok {
		tr.ResponseHeaderTimeout = time.Duration(cfg.Upstream.TimeoutSeconds) * time.Second
	}

	sch := scheduler.New(scheduler.Config{
		Pool:         p,
		Upstream:     up,
		CheckinHour:  cfg.Schedule.CheckinHour,
		RefreshHours: cfg.Schedule.RefreshHours,
		RefreshSkew:  24 * time.Hour,
	})

	h := server.NewHandler(server.Config{
		Pool:                   p,
		Upstream:               up,
		APIKey:                 cfg.APIKey,
		AuthDir:                cfg.AuthDir,
		Listen:                 cfg.Listen,
		CallbackPort:           cfg.CallbackPort,
		StateFile:              cfg.StateFile,
		UpstreamTimeoutSeconds: cfg.Upstream.TimeoutSeconds,
		PlanCooldown:           cfg.PlanCreditDur,
		SoftCooldown:           cfg.SoftRateDur,
		ErrThreshold:           cfg.Cooldown.ErrThresh,
		ErrCooldown:            cfg.ErrCooldownDur,
		DefaultModel:           cfg.DefaultModel,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go sch.Run(ctx)

	srv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           h,
		ReadHeaderTimeout: 30 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	// 第二个 http.Server：监听 CallbackPort（默认 18080），只处理 /authorize 回调。
	// 复用同一 Handler（/authorize 已在主 mux 注册）。
	// cfg.CallbackPort == "0" 时不启动（纯手动粘贴模式）。
	var cbSrv *http.Server
	if cfg.CallbackPort != "" && cfg.CallbackPort != "0" {
		cbSrv = &http.Server{
			Addr:              "127.0.0.1:" + cfg.CallbackPort,
			Handler:           h,
			ReadHeaderTimeout: 30 * time.Second,
		}
		go func() {
			<-ctx.Done()
			sc, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = cbSrv.Shutdown(sc)
		}()
		go func() {
			log.Printf("trae2api-web callback server on 127.0.0.1:%s (TRAE login /authorize)", cfg.CallbackPort)
			if err := cbSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				// 端口被占用（login.sh / 旧实例）不致命，降级为手动粘贴模式。
				log.Printf("callback server (:%s) failed: %v — web 登录降级为手动粘贴回调链接", cfg.CallbackPort, err)
			}
		}()
	}

	log.Printf("trae2api-web listening on %s (api_key=%v)", cfg.Listen, cfg.APIKey != "")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http: %v", err)
	}
	log.Printf("bye")
}
