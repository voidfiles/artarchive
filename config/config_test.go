package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppConfig(t *testing.T) {
	var testTable = []struct {
		envs   map[string]string
		output AppConfig
	}{
		{
			map[string]string{},
			AppConfig{
				Database: DatabaseConfig{Type: "postgres", DatabaseURL: defaultConfig["DATABASE_URL"]},
			}},
		{
			map[string]string{"DATABASE_URL": "postgres://remote/awesome"},
			AppConfig{Database: DatabaseConfig{Type: "postgres", DatabaseURL: "postgres://remote/awesome"}}},
	}

	for _, test := range testTable {
		defer os.Clearenv()
		for key, val := range test.envs {
			os.Setenv(key, val)
		}
		foundConfig := NewAppConfig()
		assert.Equal(t, test.output.Database, foundConfig.Database)
	}

}
