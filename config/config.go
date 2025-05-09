package config

import (
	"fmt"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Auth   AuthConfig
}

type ServerConfig struct {
	GatewayPort string `env:"GATEWAY_PORT,required"`
	UserPort    string `env:"USER_PORT,required"`
	ProductPort string `env:"PRODUCT_PORT,required"`
}

type ServiceConfig struct {
	UserServiceAddr    string `env:"USER_SERVICE_URL,required"`
	ProductServiceAddr string `env:"USER_SERVICE_URL,required"`
}

type AuthConfig struct {
	JWTSecretKey string `env:"JWT_SECRET,required"`
}

type DBConfig struct {
	Port     string `env:"DB_PORT,required"`
	Host     string `env:"DB_HOST,required"`
	Name     string `env:"DB_NAME,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config
	if err := envdecode.StrictDecode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to loading .env: %v", err)
	}
	return &cfg, nil
}
