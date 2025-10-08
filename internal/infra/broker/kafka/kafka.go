package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	wbkafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/infra"
	"github.com/sunr3d/image-processor/models"
)

var _ infra.Broker = (*KafkaBroker)(nil)

type KafkaBroker struct {
	producer *wbkafka.Producer
	consumer *wbkafka.Consumer
	topic    string
}

// New - конструктор-обертка для Kafka.
func New(brokers []string, topic, groupID string) *KafkaBroker {
	producer := wbkafka.NewProducer(brokers, topic)
	consumer := wbkafka.NewConsumer(brokers, topic, groupID)

	return &KafkaBroker{
		producer: producer,
		consumer: consumer,
		topic:    topic,
	}
}

// Publish - отправляет задачу обработки изображения в очередь Kafka.
func (kb *KafkaBroker) Publish(ctx context.Context, task *models.ProcessingTask) error {
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  2,
	}

	if err := kb.producer.SendWithRetry(ctx, strategy, []byte(task.ImageID), data); err != nil {
		return fmt.Errorf("producer.SendWithRetry: %w", err)
	}
	zlog.Logger.Info().Msgf("Задача обработки изображения отправлена в Kafka: %s", task.ImageID)

	return nil
}

// Subscribe - подписывается на очередь Kafka и выполняет обработку задач handler.
func (kb *KafkaBroker) Subscribe(ctx context.Context, handler func(ctx context.Context, task *models.ProcessingTask) error) error {
	msgChan := make(chan kafka.Message)

	strategy := retry.Strategy{
		Attempts: 3,
		Delay: time.Second,
		Backoff: 2,
	}

	kb.consumer.StartConsuming(ctx, msgChan, strategy)

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