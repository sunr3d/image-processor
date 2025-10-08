package imagesvc

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sunr3d/image-processor/mocks"
	"github.com/sunr3d/image-processor/models"
)

// UploadImage tests.
func TestImageService_UploadImage_OK(t *testing.T) {
	ctx := context.Background()
	imgStorage := mocks.NewImageStorage(t)
	metaStorage := mocks.NewMetadataStorage(t)
	publisher := mocks.NewPublisher(t)

	imgStorage.EXPECT().
		SaveOriginal(ctx, mock.AnythingOfType("string"), mock.Anything, "test.jpg").
		Return("/path/to/original", nil).
		Once()

	metaStorage.EXPECT().
		Save(ctx, mock.AnythingOfType("*models.ImageMetadata")).
		Return(nil).
		Once()

	publisher.EXPECT().
		Publish(ctx, mock.AnythingOfType("*models.ProcessingTask")).
		Return(nil).
		Once()

	svc := New(imgStorage, metaStorage, publisher)

	content := []byte("test image content")
	reader := bytes.NewReader(content)
	file := &mockMultipartFile{reader: reader, filename: "test.jpg"}

	id, err := svc.UploadImage(ctx, file, "test.jpg")

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestImageService_UploadImage_Error(t *testing.T) {
	ctx := context.Background()
	imgStorage := mocks.NewImageStorage(t)
	metaStorage := mocks.NewMetadataStorage(t)
	publisher := mocks.NewPublisher(t)

	imgStorage.EXPECT().
		SaveOriginal(ctx, mock.AnythingOfType("string"), mock.Anything, "test.jpg").
		Return("", assert.AnError).
		Once()

	svc := New(imgStorage, metaStorage, publisher)

	content := []byte("test image content")
	reader := bytes.NewReader(content)
	file := &mockMultipartFile{reader: reader, filename: "test.jpg"}

	id, err := svc.UploadImage(ctx, file, "test.jpg")

	assert.Error(t, err)
	assert.Empty(t, id)
}

// GetImage tests.
func TestImageService_GetImage_OK(t *testing.T) {
	ctx := context.Background()
	imgStorage := mocks.NewImageStorage(t)
	metaStorage := mocks.NewMetadataStorage(t)
	publisher := mocks.NewPublisher(t)

	metaStorage.EXPECT().
		Get(ctx, "test-id").
		Return(&models.ImageMetadata{ID: "test-id"}, nil).
		Once()

	imgStorage.EXPECT().
		GetPath("test-id", "original").
		Return("/path/to/image", nil).
		Once()

	svc := New(imgStorage, metaStorage, publisher)

	path, err := svc.GetImage(ctx, "test-id", "original")

	assert.NoError(t, err)
	assert.Equal(t, "/path/to/image", path)
}

// DeleteImage tests.
func TestImageService_DeleteImage_OK(t *testing.T) {
	ctx := context.Background()
	imgStorage := mocks.NewImageStorage(t)
	metaStorage := mocks.NewMetadataStorage(t)
	publisher := mocks.NewPublisher(t)

	metaStorage.EXPECT().
		Get(ctx, "test-id").
		Return(&models.ImageMetadata{ID: "test-id"}, nil).
		Once()

	imgStorage.EXPECT().
		DeleteImage(ctx, "test-id").
		Return(nil).
		Once()

	metaStorage.EXPECT().
		Delete(ctx, "test-id").
		Return(nil).
		Once()

	svc := New(imgStorage, metaStorage, publisher)

	err := svc.DeleteImage(ctx, "test-id")

	assert.NoError(t, err)
}

// Mock для multipart.File
type mockMultipartFile struct {
	reader   io.Reader
	filename string
}

func (m *mockMultipartFile) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func (m *mockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.EOF
}

func (m *mockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (m *mockMultipartFile) Size() int64 {
	return 0
}
