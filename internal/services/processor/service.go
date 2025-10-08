package processor

import (
	"fmt"

	"github.com/disintegration/imaging"
	"github.com/wb-go/wbf/zlog"
)

type ImageProcessor struct {
	thumbnailSize int
	resizeWidth   int
	watermarkText string
}

type ProcessedImages struct {
	Resized     []byte
	Thumbnail   []byte
	Watermarked []byte
}

// New - конструктор для ImageProcessor
func New(thumbSize, resizeW int, watermark string) *ImageProcessor {
	return &ImageProcessor{
		thumbnailSize: thumbSize,
		resizeWidth:   resizeW,
		watermarkText: watermark,
	}
}

// Process - обрабатывает изображение в трех вариантах: resized, thumbnail, watermarked.
func (p *ImageProcessor) Process(imagePath string) (*ProcessedImages, error) {
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

	return &ProcessedImages{
		Resized: resizedBytes,
		Thumbnail: thumbnailBytes,
		Watermarked: watermarkedBytes,
	}, nil
}
