package services

import (
	"context"
	"mime/multipart"

	"github.com/sunr3d/image-processor/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=ImageProcessor --output=../../../mocks --filename=mock_image_processor.go --with-expecter
type ImageProcessor interface {
	UploadImage(ctx context.Context, file multipart.File, filename string) (string, error)
	GetImage(ctx context.Context, id, imageType string) (string, error)
	DeleteImage(ctx context.Context, id string) error
	GetImageMeta(ctx context.Context, id string) (*models.ImageMetadata, error)
}
