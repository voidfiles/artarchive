package config

import (
	"github.com/spf13/viper"
)

var (
	defaultConfig = map[string]string{
		"DATABASE_URL":  "postgres://localhost/artarchive",
		"BUCKET":        "art.rumproarious.com",
		"VERSION":       "v2",
		"IMAGE_PATH":    "images",
		"AUTH_PASSWORD": "a-fun-password",
		"PORT":          "8000",
	}
)

type DatabaseConfig struct {
	Type        string
	DatabaseURL string
}

type AppConfig struct {
	viper        *viper.Viper
	Bucket       string
	Version      string
	ImagePath    string
	Port         string
	AuthPassword string
	Database     DatabaseConfig
}

func NewAppConfig() AppConfig {
	v := viper.New()

	for key, val := range defaultConfig {
		v.SetDefault(key, val)
		v.BindEnv(key)
	}

	return AppConfig{
		viper:        v,
		Bucket:       v.GetString("BUCKET"),
		Version:      v.GetString("VERSION"),
		ImagePath:    v.GetString("IMAGE_PATH"),
		Port:         v.GetString("PORT"),
		AuthPassword: v.GetString("AUTH_PASSWORD"),
		Database: DatabaseConfig{
			Type:        "postgres",
			DatabaseURL: v.GetString("DATABASE_URL"),
		},
	}
}
