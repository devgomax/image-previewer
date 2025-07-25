package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"

	imagepreviewer "github.com/devgomax/image-previewer/internal/app/image_previewer"
	"github.com/devgomax/image-previewer/internal/config"
	"github.com/devgomax/image-previewer/internal/logger"
	"github.com/devgomax/image-previewer/internal/pkg/lru"
	"github.com/devgomax/image-previewer/internal/pkg/resizing"
	internalhttp "github.com/devgomax/image-previewer/internal/server/http"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to get service config") //nolint:gocritic
	}

	if err = logger.ConfigureLogging(cfg.Logger); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to configure logging")
	}

	cache := lru.NewCache(cfg.LRUCacheConfig.Size)

	app := imagepreviewer.NewApp(cache, resizing.NewResizer())

	r := internalhttp.NewRouter(app,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(20*time.Second),
	)

	server := internalhttp.NewServer(cfg.HTTPConfig.GetAddr(), r)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer shutdownCancel()

		if err = server.Stop(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("failed to stop HTTP server")
		}
	}()

	if err = server.Start(ctx); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("failed to run HTTP server")
	}

	wg.Wait()
}
