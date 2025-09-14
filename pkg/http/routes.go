package http

import (
	"net/http"
	"url-shortener/pkg/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, h *handlers.Handlers) {
	huma.Post(api, "/shorten", h.Create, func(op *huma.Operation) {
		op.Description = "Создать короткую ссылку."
		op.DefaultStatus = http.StatusCreated
	})
	huma.Get(api, "/{short_id}", h.Resolve, func(op *huma.Operation) {
		op.Description = "Перенаправить по короткой ссылке (307 Redirect)."
	})
}
