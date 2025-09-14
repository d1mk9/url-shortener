package repository

import (
	"context"
	"errors"
	"url-shortener/pkg/models"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrExpired      = errors.New("expired")
	ErrLimitReached = errors.New("limit reached")
)

type Repository interface {
	Create(ctx context.Context, sl *models.ShortLink) error
	FindByShortID(ctx context.Context, shortID string) (*models.ShortLink, error)
	Resolve(ctx context.Context, shortID string) (string, error)
	Delete(ctx context.Context, shortID string) error
}
