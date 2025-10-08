package infra

import (
	"context"

	"github.com/sunr3d/image-processor/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Broker --output=../../../mocks --filename=mock_broker.go --with-expecter
type Broker interface {
	Publish(ctx context.Context, task *models.ProcessingTask) error
	Subscribe(ctx context.Context, handler func(ctx context.Context, task *models.ProcessingTask) error) error
}
