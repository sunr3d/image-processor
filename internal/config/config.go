package config

type Config struct {
	HTTPPort      string `mapstructure:"HTTP_PORT"`
	LogLevel      string `mapstructure:"LOG_LEVEL"`
	KafkaBrokers  string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic    string `mapstructure:"KAFKA_TOPIC"`
	KafkaGroup    string `mapstructure:"KAFKA_GROUP"`
	StoragePath   string `mapstructure:"STORAGE_PATH"`
	ThumbnailSize int    `mapstructure:"THUMBNAIL_SIZE"`
	ResizeWidth   int    `mapstructure:"RESIZE_WIDTH"`
	WatermarkText string `mapstructure:"WATERMARK_TEXT"`
}
