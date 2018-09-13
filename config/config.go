package config

import (
	"github.com/spf13/viper"
)

var (
	defaultConfig = map[string]string{
		"DATABASE_URL": "postgres://localhost/artarchive",
		"BUCKET":       "art.rumproarious.com",
		"VERSION":      "v2",
		"IMAGE_PATH":   "images",
	}
)

type DatabaseConfig struct {
	Type        string
	DatabaseURL string
}

type AppConfig struct {
	viper     *viper.Viper
	Bucket    string
	Version   string
	ImagePath string
	Database  DatabaseConfig
}

func NewAppConfig() AppConfig {
	v := viper.New()

	for key, val := range defaultConfig {
		v.SetDefault(key, val)
		v.BindEnv(key)
	}

	return AppConfig{
		viper:     v,
		Bucket:    v.GetString("BUCKET"),
		Version:   v.GetString("VERSION"),
		ImagePath: v.GetString("IMAGE_PATH"),
		Database: DatabaseConfig{
			Type:        "postgres",
			DatabaseURL: v.GetString("DATABASE_URL"),
		},
	}
}
