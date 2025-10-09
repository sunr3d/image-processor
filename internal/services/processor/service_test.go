package processor

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageProcessor_New(t *testing.T) {
	processor := New(200, 800)

	assert.Equal(t, 200, processor.thumbnailSize)
	assert.Equal(t, 800, processor.resizeWidth)
}

// Process tests.
func TestImageProcessor_Process_OK(t *testing.T) {
	processor := New(200, 800)
	testImagePath := createTestImage(t)
	defer os.Remove(testImagePath)

	result, err := processor.Process(testImagePath)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Resized)
	assert.NotEmpty(t, result.Thumbnail)
	assert.NotEmpty(t, result.Watermarked)
}

func TestImageProcessor_Process_FileNotFound(t *testing.T) {
	processor := New(200, 800)
	nonExistentPath := "/path/to/non/existent/image.jpg"

	result, err := processor.Process(nonExistentPath)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "imaging.Open")
}

func TestImageProcessor_Process_InvalidFile(t *testing.T) {
	processor := New(200, 800)

	tempFile := filepath.Join(t.TempDir(), "invalid.jpg")
	err := os.WriteFile(tempFile, []byte("not an image"), 0644)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	result, err := processor.Process(tempFile)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestImageProcessor_encodeToJPEG(t *testing.T) {
	img := createSimpleImage(100, 100)

	data, err := encodeToJPEG(img)

	require.NoError(t, err)
	assert.NotEmpty(t, data)

	_, err = jpeg.Decode(bytes.NewReader(data))
	assert.NoError(t, err)
}

func TestImageProcessor_addWatermark(t *testing.T) {
	processor := New(200, 800)
	img := createSimpleImage(100, 100)

	watermarked := processor.addWatermark(img)

	assert.NotNil(t, watermarked)
	assert.Equal(t, img.Bounds(), watermarked.Bounds())
}

// Helper functions
func createTestImage(t *testing.T) string {
	img := createSimpleImage(100, 100)

	tempFile := filepath.Join(t.TempDir(), "test.jpg")
	file, err := os.Create(tempFile)
	require.NoError(t, err)
	defer file.Close()

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	require.NoError(t, err)

	return tempFile
}

func createSimpleImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	return img
}
