package config

import (
	"fmt"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/zlog"
)

func GetConfig(path string) (*Config, error) {
	cfg := config.New()
	if err := cfg.Load("config.yml", ".env", ""); err != nil {
		zlog.Logger.Warn().Msgf("config.Load(): %v. Продолжаем с дефолтными значениями...", err)
	}

	cfg.SetDefault("HTTP_PORT", "8080")
	cfg.SetDefault("LOG_LEVEL", "info")
	cfg.SetDefault("KAFKA_BROKERS", "localhost:9092")
	cfg.SetDefault("KAFKA_TOPIC", "image-processing")
	cfg.SetDefault("KAFKA_GROUP", "image-processor-group")
	cfg.SetDefault("STORAGE_PATH", "./storage")
	cfg.SetDefault("METADATA_PATH", "./metadata")
	cfg.SetDefault("THUMBNAIL_SIZE", 200)
	cfg.SetDefault("RESIZE_WIDTH", 800)
	cfg.SetDefault("WATERMARK_TEXT", "© Sunr3d's Image Processor")

	var c Config
	if err := cfg.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("cfg.Unmarshal: %w", err)
	}

	return &c, nil
}
