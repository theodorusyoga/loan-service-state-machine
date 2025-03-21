package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type DatabaseType string

const (
	DatabaseTypePostgres DatabaseType = "cockroach"
)

// Config holds the service configuration
type Config struct {
	Server struct {
		Port string `yaml:"port"`
	}

	Database struct {
		Type DatabaseType `yaml:"type"`
		URL  string       `yaml:"url"`
	}
}

// Load loads configuration from YAML file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Override with environment variables if set (for deployment purposes)
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}

	if dbType := os.Getenv("DATABASE_TYPE"); dbType != "" {
		config.Database.Type = DatabaseType(dbType)
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.Database.URL = dbURL
	}

	return &config, nil
}
