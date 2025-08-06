package config


type Config struct {
	TgToken string
	
}
func Load() (*Config, error)