package filestorage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/infra"
)

var _ infra.ImageStorage = (*fileStorage)(nil)

type fileStorage struct {
	basePath string
}

// NewFileStorage - конструктор FileStorage.
func NewFileStorage(basePath string) *fileStorage {
	return &fileStorage{
		basePath: basePath,
	}
}

// SaveOriginal - сохраняет оригинал изображения и возвращает путь к нему.
func (fs *fileStorage) SaveOriginal(ctx context.Context, id string, file multipart.File, filename string) (string, error) {
	dir := filepath.Join(fs.basePath, "original", id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	path := filepath.Join(dir, filename)
	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("os.Create: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	zlog.Logger.Info().Msgf("Оригинал изображения сохранен в %s", path)

	return path, nil
}

// SaveProcessed - сохраняет обработанное изображение с указанием типа и возвращает путь к нему.
func (fs *fileStorage) SaveProcessed(ctx context.Context, id, imageType string, data []byte) (string, error) {
	dir := filepath.Join(fs.basePath, "processed", id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	filename := fmt.Sprintf("%s.jpg", imageType)
	path := filepath.Join(dir, filename)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("os.WriteFile: %w", err)
	}

	zlog.Logger.Info().Msgf("Обработанное изображение (type: %s) сохранено в %s", imageType, path)

	return path, nil
}

// GetPath - находит путь к изображению по его ID и типу.
func (fs *fileStorage) GetPath(id, imageType string) (string, error) {
	path := ""

	switch imageType {
	case "original":
		path = filepath.Join(fs.basePath, "original", id)
	case "resized", "thumbnail", "watermarked":
		filename := fmt.Sprintf("%s.jpg", imageType)
		path = filepath.Join(fs.basePath, "processed", id, filename)
	default:
		return "", fmt.Errorf("неизвестный тип изображения: %s", imageType)
	}

	return path, nil
}

// DeleteImage - удаляет изображение (оригинал и обработанные версии) по его ID.
func (fs *fileStorage) DeleteImage(ctx context.Context, id string) error {
	originalPath := filepath.Join(fs.basePath, "original", id)
	processedPath := filepath.Join(fs.basePath, "processed", id)

	if err := os.RemoveAll(originalPath); err != nil {
		zlog.Logger.Warn().Err(err).Msgf("Ошибка удаления оригинального изображения: %s", originalPath)
	}

	if err := os.RemoveAll(processedPath); err != nil {
		zlog.Logger.Warn().Err(err).Msgf("Ошибка удаления обработанного изображения: %s", processedPath)
	}

	zlog.Logger.Info().Msgf("Изображение %s успешно удалено", id)

	return nil
}
