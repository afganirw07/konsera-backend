package config

import (
	"os"
	"strings"
)

type MidtransConfig struct {
	KeyName      string
	ServerKey    string
	ClientKey    string
	IsProduction bool
	SnapURL      string
	CoreURL      string
}

func LoadMidtransConfig() MidtransConfig {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if strings.TrimSpace(serverKey) == "" {
		serverKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	}

	return MidtransConfig{
		KeyName:      os.Getenv("MIDTRANS_KEY_NAME"),
		ServerKey:    serverKey,
		ClientKey:    os.Getenv("MIDTRANS_CLIENT_KEY"),
		IsProduction: strings.EqualFold(strings.TrimSpace(os.Getenv("MIDTRANS_IS_PRODUCTION")), "true"),
		SnapURL:      os.Getenv("MIDTRANS_SNAP_URL"),
		CoreURL:      os.Getenv("MIDTRANS_CORE_URL"),
	}
}
