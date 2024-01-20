package models

import (
	"time"
)

type Post struct {
	ID          string `gorm:"primaryKey;default:cuid()"`
	Title       string
	Description string
	CreatedAt   time.Time
}
