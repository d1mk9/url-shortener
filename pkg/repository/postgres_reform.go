package repository

import (
	"context"
	"errors"
	"time"
	"url-shortener/pkg/models"

	"github.com/lib/pq"
	"gopkg.in/reform.v1"
)

type PostgresRepository struct {
	db *reform.DB
}

func NewPostgresRepository(db *reform.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, sl *models.ShortLink) error {
	if err := r.db.WithContext(ctx).Insert(sl); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) FindByShortID(ctx context.Context, shortID string) (*models.ShortLink, error) {
	sl := new(models.ShortLink)
	if err := r.db.WithContext(ctx).FindOneTo(sl, "short_id", shortID); err != nil {
		if errors.Is(err, reform.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return sl, nil
}

func (r *PostgresRepository) Resolve(ctx context.Context, shortID string) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var sl models.ShortLink
	if err := tx.WithContext(ctx).SelectOneTo(&sl, "WHERE short_id = $1 FOR UPDATE", shortID); err != nil {
		if errors.Is(err, reform.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}

	now := time.Now().UTC()

	if sl.ExpiresAt != nil && now.After(*sl.ExpiresAt) {
		return "", ErrExpired
	}

	if sl.MaxVisits > 0 && sl.Visits >= sl.MaxVisits {
		return "", ErrLimitReached
	}

	sl.Visits++
	if err := tx.WithContext(ctx).Update(&sl); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return sl.OriginalURL, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, shortID string) error {
	n, err := r.db.WithContext(ctx).DeleteFrom(models.ShortLinkTable, "WHERE short_id = $1", shortID)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
