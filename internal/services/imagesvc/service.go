package imagesvc

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/infra"
	"github.com/sunr3d/image-processor/internal/interfaces/services"
	"github.com/sunr3d/image-processor/models"
)

var _ services.ImageService = (*imageService)(nil)

type imageService struct {
	imgStorage  infra.ImageStorage
	metaStorage infra.MetadataStorage
	publisher   infra.Publisher
}

// New - конструктор imageService.
func New(imgStorage infra.ImageStorage, metaStorage infra.MetadataStorage, publisher infra.Publisher) *imageService {
	return &imageService{
		imgStorage:  imgStorage,
		metaStorage: metaStorage,
		publisher:   publisher,
	}
}

// UploadImage - загружает оригинальное изображение, сохраняет метаданные и передает задачу на обработку в брокер.
func (is *imageService) UploadImage(ctx context.Context, file multipart.File, filename string) (string, error) {
	id := uuid.New().String()

	zlog.Logger.Info().Msgf("Начало загрузки изображения: %s (ID: %s)", filename, id)

	path, err := is.imgStorage.SaveOriginal(ctx, id, file, filename)
	if err != nil {
		return "", fmt.Errorf("imgStorage.SaveOriginal: %w", err)
	}

	meta := &models.ImageMetadata{
		ID:           id,
		OriginalName: filename,
		OriginalPath: path,
		Status:       models.StatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := is.metaStorage.Save(ctx, meta); err != nil {
		return "", fmt.Errorf("metaStorage.Save: %w", err)
	}

	task := &models.ProcessingTask{
		ImageID:      id,
		OriginalPath: path,
	}

	if err := is.publisher.Publish(ctx, task); err != nil {
		return "", fmt.Errorf("broker.Publish: %w", err)
	}

	zlog.Logger.Info().Msgf("Изображение %s успешно загружено и передано на обработку", id)

	return id, nil
}

// GetImage - получает путь к изображению по его ID и типу.
func (is *imageService) GetImage(ctx context.Context, id, imageType string) (string, error) {
	meta, err := is.metaStorage.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("metaStorage.Get: %w", err)
	}

	if imageType != "original" && meta.Status != models.StatusCompleted {
		return "", fmt.Errorf("изображение еще не обработано, статус: %s", meta.Status)
	}

	path, err := is.imgStorage.GetPath(id, imageType)
	if err != nil {
		return "", fmt.Errorf("imgStorage.GetPath: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("изображение не найдено: %s", path)
	}

	return path, nil
}

// DeleteImage - удаляет изображение по его ID.
func (is *imageService) DeleteImage(ctx context.Context, id string) error {
	if _, err := is.metaStorage.Get(ctx, id); err != nil {
		return fmt.Errorf("metaStorage.Get: %w", err)
	}

	if err := is.metaStorage.Delete(ctx, id); err != nil {
		return fmt.Errorf("metaStorage.Delete: %w", err)
	}

	if err := is.imgStorage.DeleteImage(ctx, id); err != nil {
		zlog.Logger.Warn().Err(err).Msgf("Ошибка удаления файлов для: %s", id)
		return nil
	}

	zlog.Logger.Info().Msgf("Изображение %s и его метаданные успешно удалены", id)

	return nil
}