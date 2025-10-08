package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"

	"github.com/disintegration/imaging"
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/services"
	"github.com/sunr3d/image-processor/models"
)

var _ services.ImageProcessor = (*imageProcessor)(nil)

type imageProcessor struct {
	thumbnailSize int
	resizeWidth   int
	watermarkText string
}

// New - конструктор для ImageProcessor
func New(thumbSize, resizeW int, watermark string) *imageProcessor {
	return &imageProcessor{
		thumbnailSize: thumbSize,
		resizeWidth:   resizeW,
		watermarkText: watermark,
	}
}

// Process - обрабатывает изображение в трех вариантах: resized, thumbnail, watermarked.
func (p *imageProcessor) Process(imagePath string) (*models.ProcessedImages, error) {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("imaging.Open: %w", err)
	}

	zlog.Logger.Info().Msgf("Начало обработки изображения: %s", imagePath)

	// 1. Resize
	resized := imaging.Resize(img, p.resizeWidth, 0, imaging.Lanczos)
	resizedBytes, err := encodeToJPEG(resized)
	if err != nil {
		return nil, fmt.Errorf("encodeToJPEG Resized: %w", err)
	}

	// 2. Thumbnail
	thumbnail := imaging.Fill(img, p.thumbnailSize, p.thumbnailSize, imaging.Center, imaging.Lanczos)
	thumbnailBytes, err := encodeToJPEG(thumbnail)
	if err != nil {
		return nil, fmt.Errorf("encodeToJPEG Thumbnail: %w", err)
	}

	// 3. Watermark
	watermarked := p.addWatermark(img)
	watermarkedBytes, err := encodeToJPEG(watermarked)
	if err != nil {
		return nil, fmt.Errorf("encodeToJPEG Watermarked: %w", err)
	}

	zlog.Logger.Info().Msgf("Обработка изображения %s завершена", imagePath)

	return &models.ProcessedImages{
		Resized:     resizedBytes,
		Thumbnail:   thumbnailBytes,
		Watermarked: watermarkedBytes,
	}, nil
}

// Helpers
func encodeToJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	opts := &jpeg.Options{Quality: 90}

	if err := jpeg.Encode(&buf, img, opts); err != nil {
		return nil, fmt.Errorf("jpeg.Encode: %w", err)
	}

	return buf.Bytes(), nil
}

func (p *imageProcessor) addWatermark(img image.Image) image.Image {
	overlay := image.NewRGBA(img.Bounds())

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if (x+y)%100 < 20 {
				overlay.Set(x, y, color.RGBA{255, 255, 255, 50})
			}
		}
	}

	return imaging.Overlay(img, overlay, image.Pt(0, 0), 0.3)
}
