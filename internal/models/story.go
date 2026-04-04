package models

import "time"

type Size string

const (
	SizeSmall Size = "small"
	SizeLarge Size = "large"
)

type Story struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	CoverImage  string    `json:"cover_image"`
	Author      string    `json:"author"`
	Content     string    `json:"content"`
	AIGenerated bool      `json:"ai_generated"`
	Size        Size      `json:"size"`
	Views       int64     `json:"views"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
