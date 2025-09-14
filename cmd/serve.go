package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"url-shortener/pkg/config"
	httpapi "url-shortener/pkg/http"
	"url-shortener/pkg/repository"
	"url-shortener/pkg/service"
	"url-shortener/pkg/storage"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запустить HTTP-сервер",

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.MustLoad()
		dbs := storage.MustInitPostgres(cfg.PostgresDSN())
		defer func() {
			if err := dbs.SQL.Close(); err != nil {
				log.Printf("db close error: %v", err)
			}
		}()
		repo := repository.NewPostgresRepository(dbs.Reform)
		gen := service.NewUUIDGen(8)
		svc := service.NewShortLinkService(repo, gen)
		srv := httpapi.NewServer(cfg, svc, cfg.BaseURL)
		addr := ":8080"

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		errCh := make(chan error, 1)
		go func() { errCh <- srv.Run(addr) }()

		select {
		case <-ctx.Done():
			log.Println("shutdown signal received")
		case err := <-errCh:
			if err != nil {
				log.Printf("server run error: %v", err)
				return err
			}
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown error: %v", err)
		}
		log.Println("server stopped gracefully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
