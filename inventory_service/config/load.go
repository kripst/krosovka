package config

import "fmt"

type StorageConfig struct {
	Host    string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
	Port    string `yaml:"port" env:"PG_PORT" env-default:"5432"`
	User    string `yaml:"user" env:"PG_USER" env-default:"postgres"`
	Pass    string `yaml:"pass" env:"PG_PASSWORD" env-default:"password"`
	DBName  string `yaml:"dbname" env:"PG_DBNAME" env-default:"database"`
	SSLMode string `yaml:"sslmode" env:"PG_SSLMODE" env-default:"disable"`
	PoolMax int   `yaml:"pool_max" env:"PG_POOL_MAX" env-default:"10"`
}

func (c *StorageConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

type Config struct {
	StorageConfig *StorageConfig
}

func Load() (*Config, error) {
	return &Config{

	}, nil
}