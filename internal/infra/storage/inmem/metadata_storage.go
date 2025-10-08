package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/infra"
	"github.com/sunr3d/image-processor/models"
)

const (
	metadataCacheSize = 128
)

var _ infra.MetadataStorage = (*MetadataStorage)(nil)

type MetadataStorage struct {
	metadata map[string]*models.ImageMetadata
	mu       sync.RWMutex
}

// NewMetadataStorage - конструктор MetadataStorage.
func NewMetadataStorage() *MetadataStorage {
	return &MetadataStorage{
		metadata: make(map[string]*models.ImageMetadata, metadataCacheSize),
	}
}

// Save - сохраняет метаданные изображения.
func (ms *MetadataStorage) Save(ctx context.Context, meta *models.ImageMetadata) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.metadata[meta.ID] = meta

	zlog.Logger.Info().Msgf("Метаданные изображения сохранены: %s", meta.ID)

	return nil
}

// Get - получает метаданные изображения по ID.
func (ms *MetadataStorage) Get(ctx context.Context, id string) (*models.ImageMetadata, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	meta, exists := ms.metadata[id]
	if !exists {
		return nil, fmt.Errorf("метаданные изображения не найдены: %s", id)
	}

	return meta, nil
}

// Update - обновляет метаданные изображения.
func (ms *MetadataStorage) Update(ctx context.Context, meta *models.ImageMetadata) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.metadata[meta.ID]; !exists {
		return fmt.Errorf("метаданные изображения не найдены: %s", meta.ID)
	}

	ms.metadata[meta.ID] = meta
	zlog.Logger.Info().Msgf("Метаданные изображения обновлены (%s): %s", meta.Status, meta.ID)

	return nil
}

// Delete - удаляет метаданные изображения по ID.
func (ms *MetadataStorage) Delete(ctx context.Context, id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.metadata, id)
	zlog.Logger.Info().Msgf("Метаданные изображения удалены: %s", id)

	return nil
}