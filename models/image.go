package models

import "time"

type ImageStatus string

const (
	StatusPending    ImageStatus = "pending"
	StatusProcessing ImageStatus = "processing"
	StatusCompleted  ImageStatus = "completed"
	StatusFailed     ImageStatus = "failed"
)

type ImageMetadata struct {
	ID              string
	OriginalName    string
	OriginalPath    string
	ResizedPath     string
	ThumbnailPath   string
	WatermarkedPath string
	Status          ImageStatus
	ErrorMessage    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ProcessedImages struct {
	Resized     []byte
	Thumbnail   []byte
	Watermarked []byte
}
