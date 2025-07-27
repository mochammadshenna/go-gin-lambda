package main

import (
	"ai-service/cmd/config"
	"ai-service/internal/app/middleware"
	"ai-service/internal/outbound"
	"ai-service/internal/routes"
	"ai-service/internal/util/authentication"
	"ai-service/internal/util/logger"
	"ai-service/internal/util/template"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// @title							AI Service API
// @version							1.0
// @description						API Doc for AI Service

// @BasePath						/api
// @securityDefinitions.apikey 		BearerAuth
// @in 								header
// @name 							Authorization
func main() {
	var (
		env                         = os.Getenv("ENV")
		port                        = os.Getenv("PORT")
		sentryUrl                   = os.Getenv("SENTRY_DSN")
		jwtSecret                   = os.Getenv("JWT_SECRET")
		signalChan chan (os.Signal) = make(chan os.Signal, 1)
	)

	// Set default port if not provided
	if port == "" {
		port = "8080"
	}

	// Initialize core components
	logger.Init(sentryUrl)
	authentication.Init(jwtSecret)
	template.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize AI manager
	aiManager := outbound.NewManager(cfg)

	startBootTime := time.Now()
	router := routes.NewRouters(aiManager)

	if env == "prod" {
		fmt.Println("running production mode")
		log.Fatal(http.ListenAndServe(":"+port, middleware.MultipleMiddleware(middleware.NewHttpMiddleware(router))))
	} else {
		fmt.Printf("running on localhost:%s \n", port)
		server := &http.Server{
			Addr:    "127.0.0.1:" + port,
			Handler: middleware.MultipleMiddleware(middleware.NewHttpMiddleware(router)),
		}

		_, cancelSubs := context.WithCancel(context.Background())

		// SIGINT handles Ctrl+C locally.
		// SIGTERM handles termination signal.
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		// Start HTTP server.
		go func() {
			logger.Info(context.Background(), "apps started", port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("http server error = %v", err)
			}
		}()
		logger.Info(context.Background(), "server booting time ", time.Since(startBootTime).Milliseconds(), "ms")

		// Receive output from signalChan.
		sig := <-signalChan
		logger.Logger.Infof("%s signal caught", sig)

		// Timeout if waiting for connections to return idle.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cancelSubs()

		// Gracefully shutdown the server by waiting on existing requests (except websockets).
		if err := server.Shutdown(ctx); err != nil {
			logger.Logger.Errorf("server shutdown failed: %+v", err)
		}

		logger.Logger.Info("server exited")

		// pprof goroutine
		isEnablePprof, _ := strconv.ParseBool(os.Getenv("ENABLE_PPROF"))
		if isEnablePprof {
			fmt.Println("PProf enabled")
		}
	}
}
