package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vo1dFl0w/marketplace-parser-service/internal/adapters/browser/chromium"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/adapters/parsers"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/config"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
	ht "github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http/httpgen"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/usecase"
	"github.com/vo1dFl0w/marketplace-parser-service/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Println(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	
	loggerCfg := logger.NewLoggerConfig(cfg.Server.Env, cfg.Options.LoggerTimeFormat)
	logger := logger.LoadLogger(loggerCfg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)

	time.Sleep(time.Second * 5)

	chromiumRepo := chromium.NewChromiumRepository(cfg)
	browser := chromium.NewBrowser(chromiumRepo)

	wb := parsers.NewWildberriesParser(cfg, logger, browser.Chromium())
	oz := parsers.NewOzonParser(cfg, logger, browser.Chromium())

	searchSvc := usecase.NewSearchService([]repository.SearchRepository{oz, wb})

	handler := ht.NewHandler(logger, searchSvc, cfg.Server.RequestTimeout)

	srv, err := httpgen.NewServer(handler)
	if err != nil {
		return fmt.Errorf("new server: %w", err)
	}

	// CORS middleware должен быть первым в цепочке
	withMiddlewares := handler.CORSMiddleware(handler.RequestTimeoutMiddleware(handler.LoggerMiddleware(srv)))

	httpServer := &http.Server{
		Addr:    cfg.Server.HTTPAddr,
		Handler: withMiddlewares,
	}

	go func() {
		logger.Info("server started", "host", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	select {
	case e := <-serverErr:
		return fmt.Errorf("server error: %w", e)
	case s := <-sig:
		logger.Info("initialization gracefull shutdown", "signal", s)
		ctxShutdown, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := httpServer.Shutdown(ctxShutdown); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}

		logger.Info("server gracefully stopped")
		return nil
	}
}
