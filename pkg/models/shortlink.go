package models

import (
	"time"

	"github.com/google/uuid"
)

//go:generate reform
//reform:short_links
type ShortLink struct {
	ID          uuid.UUID  `reform:"id,pk"`
	ShortID     string     `reform:"short_id"`
	OriginalURL string     `reform:"original_url"`
	CreatedAt   time.Time  `reform:"created_at"`
	ExpiresAt   *time.Time `reform:"expires_at"`
	MaxVisits   int64      `reform:"max_visits"`
	Visits      int64      `reform:"visits"`
}
