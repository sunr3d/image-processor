package worker

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sunr3d/image-processor/mocks"
	"github.com/sunr3d/image-processor/models"
)

func TestWorker_New(t *testing.T) {
	mockProcessor := mocks.NewImageProcessor(t)
	mockImgStorage := mocks.NewImageStorage(t)
	mockMetaStorage := mocks.NewMetadataStorage(t)
	mockSubscriber := mocks.NewSubscriber(t)

	worker := New(mockProcessor, mockImgStorage, mockMetaStorage, mockSubscriber)

	assert.NotNil(t, worker)
}

// ProcessTask tests.
func TestWorker_ProcessTask_OK(t *testing.T) {
	ctx := context.Background()

	mockProcessor := mocks.NewImageProcessor(t)
	mockImgStorage := mocks.NewImageStorage(t)
	mockMetaStorage := mocks.NewMetadataStorage(t)
	mockSubscriber := mocks.NewSubscriber(t)

	mockProcessor.EXPECT().
		Process("/path/to/original").
		Return(&models.ProcessedImages{
			Resized:     []byte("resized data"),
			Thumbnail:   []byte("thumbnail data"),
			Watermarked: []byte("watermarked data"),
		}, nil).
		Once()

	mockMetaStorage.EXPECT().
		Get(ctx, "test-id").
		Return(&models.ImageMetadata{ID: "test-id"}, nil).
		Once()

	mockMetaStorage.EXPECT().
		Update(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(nil).
		Once()

	mockImgStorage.EXPECT().
		SaveProcessed(ctx, "test-id", "resized", []byte("resized data")).
		Return("/path/to/resized", nil).
		Once()

	mockImgStorage.EXPECT().
		SaveProcessed(ctx, "test-id", "thumbnail", []byte("thumbnail data")).
		Return("/path/to/thumbnail", nil).
		Once()

	mockImgStorage.EXPECT().
		SaveProcessed(ctx, "test-id", "watermarked", []byte("watermarked data")).
		Return("/path/to/watermarked", nil).
		Once()

	mockMetaStorage.EXPECT().
		Update(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(nil).
		Once()

	worker := New(mockProcessor, mockImgStorage, mockMetaStorage, mockSubscriber)

	task := &models.ProcessingTask{
		ImageID:      "test-id",
		OriginalPath: "/path/to/original",
	}

	err := worker.processTask(ctx, task)

	assert.NoError(t, err)
}

func TestWorker_ProcessTask_Error(t *testing.T) {
	ctx := context.Background()

	mockProcessor := mocks.NewImageProcessor(t)
	mockImgStorage := mocks.NewImageStorage(t)
	mockMetaStorage := mocks.NewMetadataStorage(t)
	mockSubscriber := mocks.NewSubscriber(t)

	mockProcessor.EXPECT().
		Process("/path/to/original").
		Return(nil, errors.New("processing failed")).
		Once()

	mockMetaStorage.EXPECT().
		Get(ctx, "test-id").
		Return(&models.ImageMetadata{ID: "test-id"}, nil).
		Once()

	mockMetaStorage.EXPECT().
		Update(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(nil).
		Once()

	mockMetaStorage.EXPECT().
		Update(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(nil).
		Once()

	worker := New(mockProcessor, mockImgStorage, mockMetaStorage, mockSubscriber)

	task := &models.ProcessingTask{
		ImageID:      "test-id",
		OriginalPath: "/path/to/original",
	}

	err := worker.processTask(ctx, task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "processing failed")
}

func TestWorker_ProcessTask_SaveImagesError(t *testing.T) {
	ctx := context.Background()

	mockProcessor := mocks.NewImageProcessor(t)
	mockImgStorage := mocks.NewImageStorage(t)
	mockMetaStorage := mocks.NewMetadataStorage(t)
	mockSubscriber := mocks.NewSubscriber(t)

	mockProcessor.EXPECT().
		Process("/path/to/original").
		Return(&models.ProcessedImages{
			Resized:     []byte("resized data"),
			Thumbnail:   []byte("thumbnail data"),
			Watermarked: []byte("watermarked data"),
		}, nil).
		Once()

	mockMetaStorage.EXPECT().
		Get(ctx, "test-id").
		Return(&models.ImageMetadata{ID: "test-id"}, nil).
		Once()

	mockMetaStorage.EXPECT().
		Update(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(nil).
		Once()

	mockImgStorage.EXPECT().
		SaveProcessed(ctx, "test-id", "resized", []byte("resized data")).
		Return("", errors.New("save failed")).
		Once()

	mockMetaStorage.EXPECT().
		Update(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(nil).
		Once()

	worker := New(mockProcessor, mockImgStorage, mockMetaStorage, mockSubscriber)

	task := &models.ProcessingTask{
		ImageID:      "test-id",
		OriginalPath: "/path/to/original",
	}

	err := worker.processTask(ctx, task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save failed")
}

func TestWorker_ProcessTask_UpdateMetadataError(t *testing.T) {
	ctx := context.Background()

	mockProcessor := mocks.NewImageProcessor(t)
	mockImgStorage := mocks.NewImageStorage(t)
	mockMetaStorage := mocks.NewMetadataStorage(t)
	mockSubscriber := mocks.NewSubscriber(t)

	mockProcessor.EXPECT().
		Process("/path/to/original").
		Return(&models.ProcessedImages{
			Resized:     []byte("resized data"),
			Thumbnail:   []byte("thumbnail data"),
			Watermarked: []byte("watermarked data"),
		}, nil).
		Once()

	mockMetaStorage.EXPECT().
		Get(ctx, "test-id").
		Return(&models.ImageMetadata{ID: "test-id"}, nil).
		Once()

	mockMetaStorage.EXPECT().
		Update(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(nil).
		Once()

	mockImgStorage.EXPECT().
		SaveProcessed(ctx, "test-id", "resized", []byte("resized data")).
		Return("/path/to/resized", nil).
		Once()

	mockImgStorage.EXPECT().
		SaveProcessed(ctx, "test-id", "thumbnail", []byte("thumbnail data")).
		Return("/path/to/thumbnail", nil).
		Once()

	mockImgStorage.EXPECT().
		SaveProcessed(ctx, "test-id", "watermarked", []byte("watermarked data")).
		Return("/path/to/watermarked", nil).
		Once()

	mockMetaStorage.EXPECT().
		Update(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(errors.New("update failed")).
		Once()

	worker := New(mockProcessor, mockImgStorage, mockMetaStorage, mockSubscriber)

	task := &models.ProcessingTask{
		ImageID:      "test-id",
		OriginalPath: "/path/to/original",
	}

	err := worker.processTask(ctx, task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
}
