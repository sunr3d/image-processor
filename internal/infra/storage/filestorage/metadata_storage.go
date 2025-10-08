package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/infra"
	"github.com/sunr3d/image-processor/models"
)

var _ infra.MetadataStorage = (*metadataStorage)(nil)

type metadataStorage struct {
	basePath string
	mu       sync.RWMutex
}

// NewMetadataStorage - конструктор MetadataStorage.
func NewMetadataStorage(basePath string) *metadataStorage {
	return &metadataStorage{
		basePath: basePath,
	}
}

// Save - сохраняет метаданные изображения.
func (ms *metadataStorage) Save(ctx context.Context, meta *models.ImageMetadata) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if err := os.MkdirAll(ms.basePath, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	path := filepath.Join(ms.basePath, fmt.Sprintf("%s.json", meta.ID))
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}

	zlog.Logger.Info().Msgf("Метаданные изображения сохранены: %s", path)

	return nil
}

// Get - получает метаданные изображения по ID.
func (ms *metadataStorage) Get(ctx context.Context, id string) (*models.ImageMetadata, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	path := filepath.Join(ms.basePath, fmt.Sprintf("%s.json", id))
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("метаданные изображения не найдены: %s", id)
		}
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	var meta models.ImageMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &meta, nil
}

// Update - обновляет метаданные изображения.
func (ms *metadataStorage) Update(ctx context.Context, meta *models.ImageMetadata) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	path := filepath.Join(ms.basePath, fmt.Sprintf("%s.json", meta.ID))

	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}

	zlog.Logger.Info().Msgf("Метаданные изображения обновлены (%s): %s", meta.ID, path)

	return nil
}

// Delete - удаляет метаданные изображения по ID.
func (ms *metadataStorage) Delete(ctx context.Context, id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	path := filepath.Join(ms.basePath, fmt.Sprintf("%s.json", id))
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("метаданные изображения не найдены: %s", id)
		}
		return fmt.Errorf("os.Remove: %w", err)
	}

	zlog.Logger.Info().Msgf("Метаданные изображения удалены (%s): %s", id, path)

	return nil
}
