package infra

import (
	"context"

	"github.com/sunr3d/image-processor/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Subscriber --output=../../../mocks --filename=mock_subscriber.go --with-expecter
type Subscriber interface {
	Subscribe(ctx context.Context, handler func(ctx context.Context, task *models.ProcessingTask) error) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Publisher --output=../../../mocks --filename=mock_publisher.go --with-expecter
type Publisher interface {
	Publish(ctx context.Context, task *models.ProcessingTask) error
}
