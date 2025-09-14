package http

import (
	"context"
	"time"

	"url-shortener/pkg/config"
	"url-shortener/pkg/handlers"
	"url-shortener/pkg/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app     *fiber.App
	api     huma.API
	service service.ShortLink
	cfg     *config.Config
}

func NewServer(cfg *config.Config, svc service.ShortLink, baseURL string) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	hcfg := huma.DefaultConfig("application/json", "utf-8")
	hcfg.Info.Title = "Url Shortener API"
	hcfg.Info.Version = "1.0.0"

	api := humafiber.New(app, hcfg)
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	h := handlers.NewShortLinkHandler(svc, baseURL)
	RegisterRoutes(api, h)

	return &Server{app: app, api: api, service: svc, cfg: cfg}
}

func (s *Server) App() *fiber.App { return s.app }

func (s *Server) Run(addr string) error { return s.app.Listen(addr) }

func (s *Server) Shutdown(ctx context.Context) error { return s.app.ShutdownWithContext(ctx) }
