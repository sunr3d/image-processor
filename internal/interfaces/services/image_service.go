package services

import (
	"context"
	"mime/multipart"

)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=ImageService --output=../../../mocks --filename=mock_image_service.go --with-expecter
type ImageService interface {
	UploadImage(ctx context.Context, file multipart.File, filename string) (string, error)
	GetImage(ctx context.Context, id, imageType string) (string, error)
	DeleteImage(ctx context.Context, id string) error
}
