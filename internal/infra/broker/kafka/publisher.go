package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	wbkafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/infra"
	"github.com/sunr3d/image-processor/models"
)

var _ infra.Publisher = (*publisher)(nil)

type publisher struct {
	producer *wbkafka.Producer
	topic    string
}

// NewPublisher - конструктор publisher.
func NewPublisher(brokers []string, topic string) *publisher {
	producer := wbkafka.NewProducer(brokers, topic)

	return &publisher{
		producer: producer,
		topic:    topic,
	}
}

// Publish - отправляет задачу обработки изображения в очередь Kafka.
func (p *publisher) Publish(ctx context.Context, task *models.ProcessingTask) error {
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  2,
	}

	if err := p.producer.SendWithRetry(ctx, strategy, []byte(task.ImageID), data); err != nil {
		return fmt.Errorf("producer.SendWithRetry: %w", err)
	}
	zlog.Logger.Info().Msgf("Задача обработки изображения отправлена в Kafka: %s", task.ImageID)

	return nil
}

func (p *publisher) Close() error {
	if err := p.producer.Close(); err != nil {
		zlog.Logger.Warn().Err(err).Msg("producer.Close")
		return err
	}

	return nil
}
