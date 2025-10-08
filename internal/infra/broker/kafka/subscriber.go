package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	wbkafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/infra"
	"github.com/sunr3d/image-processor/models"
)

var _ infra.Subscriber = (*subscriber)(nil)

type subscriber struct {
	consumer *wbkafka.Consumer
	topic    string
}

// NewSubscriber - конструктор subscriber.
func NewSubscriber(brokers []string, topic, groupID string) *subscriber {
	consumer := wbkafka.NewConsumer(brokers, topic, groupID)

	return &subscriber{
		consumer: consumer,
		topic:    topic,
	}
}

// Subscribe - подписывается на очередь Kafka и выполняет обработку задач handler.
func (s *subscriber) Subscribe(ctx context.Context, handler func(ctx context.Context, task *models.ProcessingTask) error) error {
	msgChan := make(chan kafka.Message)

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  2,
	}

	s.consumer.StartConsuming(ctx, msgChan, strategy)

	for {
		select {
		case <-ctx.Done():
			zlog.Logger.Info().Msg("Получен сигнал завершения контекста, остановка подписки на Kafka")
			return nil
		case msg := <-msgChan:
			var task models.ProcessingTask
			if err := json.Unmarshal(msg.Value, &task); err != nil {
				zlog.Logger.Error().Err(err).Msgf("Ошибка при разборе задачи обработки изображения из Kafka: %s", msg.Value)
				continue
			}

			zlog.Logger.Info().Msgf("Получена задача обработки изображения из Kafka: %s", task.ImageID)

			if err := handler(ctx, &task); err != nil {
				zlog.Logger.Error().Err(err).Msgf("Ошибка при обработке задачи обработки изображения из Kafka: %s", task.ImageID)
			}
		}
	}
}

func (s *subscriber) Close() error {
	if err := s.consumer.Close(); err != nil {
		zlog.Logger.Warn().Err(err).Msg("consumer.Close")
		return err
	}

	return nil
}
