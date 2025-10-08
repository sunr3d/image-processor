package services

import "github.com/sunr3d/image-processor/models"

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=ImageProcessor --output=../../../mocks --filename=mock_image_processor.go --with-expecter
type ImageProcessor interface {
	Process(imagePath string) (*models.ProcessedImages, error)
}
