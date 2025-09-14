package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"url-shortener/pkg/service"

	"github.com/danielgtaylor/huma/v2"
)

type ShortenRequest struct {
	OriginalURL string `required:"true" format:"uri"`
	MaxVisits   *int64 `json:",omitempty" minimum:"0"`
	Expiry      string `json:",omitempty"`
}

type ShortenInput struct {
	Body ShortenRequest
}

type ShortenResponse struct {
	ShortURL  string
	ShortID   string
	ExpiresAt *time.Time `json:",omitempty"`
	MaxVisits *int64
}

type ShortenOutput struct {
	Body ShortenResponse
}

type ResolveInput struct {
	ShortID string `path:"short_id"`
}

type Handlers struct {
	svc     service.ShortLink
	baseURL string
}

func NewShortLinkHandler(svc service.ShortLink, baseURL string) *Handlers {
	return &Handlers{svc: svc, baseURL: baseURL}
}

func (h *Handlers) Create(ctx context.Context, in *ShortenInput) (*ShortenOutput, error) {
	u, _ := url.Parse(in.Body.OriginalURL)
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, huma.Error400BadRequest("original_url must use http or https")
	}

	var ttl *time.Duration
	if s := strings.TrimSpace(in.Body.Expiry); s != "" {
		if d, err := time.ParseDuration(s); err == nil && d > 0 {
			ttl = &d
		} else {
			return nil, huma.Error400BadRequest("invalid expiry duration")
		}
	}

	mv := in.Body.MaxVisits
	if mv != nil && *mv == 0 {
		mv = nil
	}

	sl, err := h.svc.CreateShortLink(ctx, u.String(), ttl, mv)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrConflict):
			return nil, huma.Error409Conflict("")
		}
	}

	return &ShortenOutput{
		Body: ShortenResponse{
			ShortURL:  h.baseURL + "/" + sl.ShortID,
			ShortID:   sl.ShortID,
			ExpiresAt: sl.ExpiresAt,
			MaxVisits: sl.MaxVisits,
		},
	}, nil
}

type redirectOut struct {
	Status  int
	Headers map[string]string
}

func (h *Handlers) Resolve(ctx context.Context, in *ResolveInput) (*redirectOut, error) {
	orig, err := h.svc.Resolve(ctx, in.ShortID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return nil, huma.Error404NotFound("")
		case errors.Is(err, service.ErrExpired), errors.Is(err, service.ErrLimitReached):
			return nil, huma.Error410Gone("")
		}
	}

	return &redirectOut{
		Status: http.StatusTemporaryRedirect,
		Headers: map[string]string{
			"Location":      orig,
			"Cache-Control": "no-store",
		},
	}, nil
}
