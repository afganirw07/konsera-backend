package config

import (
	"os"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins []string
}

func LoadCORSConfig() CORSConfig {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" {
		return CORSConfig{}
	}

	origins := make([]string, 0)
	for _, origin := range strings.Split(raw, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return CORSConfig{AllowedOrigins: origins}
}
