package infra

import (
	"context"
	"mime/multipart"

	"github.com/sunr3d/image-processor/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=ImageStorage --output=../../../mocks --filename=mock_image_storage.go --with-expecter
type ImageStorage interface {
	SaveOriginal(ctx context.Context, id string, file multipart.File, filename string) (string, error)
	SaveProcessed(ctx context.Context, id, imageType string, data []byte) (string, error)
	GetPath(id, imageType string) (string, error)
	DeleteImage(ctx context.Context, id string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=MetadataStorage --output=../../../mocks --filename=mock_metadata_storage.go --with-expecter
type MetadataStorage interface {
	Save(ctx context.Context, meta *models.ImageMetadata) error
	Get(ctx context.Context, id string) (*models.ImageMetadata, error)
	Update(ctx context.Context, meta *models.ImageMetadata) error
	Delete(ctx context.Context, id string) error
}
