package httphandlers

func validateContentType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
	}

	for _, validType := range validTypes {
		if validType == contentType {
			return true
		}
	}

	return false
}

func validateImgType(imageType string) bool {
	validTypes := []string{
		"original",
		"resized",
		"thumbnail",
		"watermarked",
	}

	for _, validType := range validTypes {
		if validType == imageType {
			return true
		}
	}

	return false
}
