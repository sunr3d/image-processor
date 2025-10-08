package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

func encodeToJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	opts := &jpeg.Options{Quality: 90}

	if err := jpeg.Encode(&buf, img, opts); err != nil {
		return nil, fmt.Errorf("jpeg.Encode: %w", err)
	}

	return buf.Bytes(), nil
}

func (p *ImageProcessor) addWatermark(img image.Image) image.Image {
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
