package config

import "os"

type ImageUploadConfig struct {
	Path   string
	APIKey string
}

func LoadImageUploadConfig() ImageUploadConfig {
	return ImageUploadConfig{
		Path:   os.Getenv("IMAGE_UPLOAD_PATH"),
		APIKey: os.Getenv("IMAGE_UPLOAD_API_KEY"),
	}
}
