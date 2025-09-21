package service

import (
	"context"
	"errors"
	"time"
	"url-shortener/pkg/models"
	"url-shortener/pkg/repository"
)

var (
	ErrConflict     = errors.New("short id conflict")
	ErrNotFound     = errors.New("not found")
	ErrExpired      = errors.New("expired")
	ErrLimitReached = errors.New("max visits reached")
)

type ShortLink interface {
	CreateShortLink(ctx context.Context, originalURL string, ttl *time.Duration, maxVisits int64) (*models.ShortLink, error)
	Resolve(ctx context.Context, shortID string) (string, error)
	Delete(ctx context.Context, shortID string) error
}

type service struct {
	repo  repository.Repository
	cache *InMemoryCache
	idGen IDGenerator
}

func NewShortLinkService(repo repository.Repository, gen IDGenerator) *service {
	return &service{
		repo:  repo,
		cache: newInMemoryCache(),
		idGen: gen,
	}
}

func (s *service) CreateShortLink(ctx context.Context, originalURL string, ttl *time.Duration, maxVisits int64) (*models.ShortLink, error) {
	now := time.Now().UTC()

	id, err := s.idGen.NewID()
	if err != nil {
		return nil, err
	}

	var expiresAt *time.Time
	if ttl != nil {
		t := now.Add(*ttl)
		expiresAt = &t
	}

	sl := &models.ShortLink{
		ShortID:     id,
		OriginalURL: originalURL,
		CreatedAt:   now,
		ExpiresAt:   expiresAt,
		MaxVisits:   maxVisits,
		Visits:      0,
	}

	if err := s.repo.Create(ctx, sl); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, ErrConflict
		}
		return nil, err
	}

	s.cache.Set(sl.ShortID, CacheEntry{
		OriginalURL: sl.OriginalURL,
		ExpiresAt:   sl.ExpiresAt,
		MaxVisits:   sl.MaxVisits,
	})

	return sl, nil
}

func (s *service) Resolve(ctx context.Context, shortID string) (string, error) {
	if ce, ok := s.cache.Get(shortID); ok && ce.ExpiresAt != nil {
		if time.Now().After(*ce.ExpiresAt) {
			return "", ErrExpired
		}
	}
	orig, err := s.repo.Resolve(ctx, shortID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return "", ErrNotFound
		case errors.Is(err, repository.ErrExpired):
			return "", ErrExpired
		case errors.Is(err, repository.ErrLimitReached):
			return "", ErrLimitReached
		default:
			return "", err
		}
	}

	return orig, nil
}

func (s *service) Delete(ctx context.Context, shortID string) error {
	if err := s.repo.Delete(ctx, shortID); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return ErrNotFound
		default:
			return err
		}
	}

	s.cache.Delete(shortID)

	return nil
}
