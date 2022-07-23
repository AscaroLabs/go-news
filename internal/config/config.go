package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	db_host     string
	db_port     string
	db_name     string
	db_user     string
	db_password string
	grpc_host   string
	grpc_port   string
	rest_host   string
	rest_port   string
	main_dir    string
	jwt_secret  string
	jwt_ttl     string
	refresh_ttl string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		db_host:     os.Getenv("DB_HOST_ADDR"),
		db_port:     os.Getenv("DB_HOST_PORT"),
		db_name:     os.Getenv("DB_NAME"),
		db_user:     os.Getenv("DB_USERNAME"),
		db_password: os.Getenv("DB_PASSWORD"),
		grpc_host:   os.Getenv("GRPC_HOST"),
		grpc_port:   os.Getenv("GRPC_PORT"),
		rest_host:   os.Getenv("REST_HOST"),
		rest_port:   os.Getenv("REST_PORT"),
		main_dir:    os.Getenv("MAIN_DIR"),
		jwt_secret:  os.Getenv("JWT_SECRET"),
		jwt_ttl:     os.Getenv("JWT_TTL"),
		refresh_ttl: os.Getenv("REFRESH_TTL"),
	}
}

func (cfg *Config) GetRefreshTTL() string {
	return cfg.refresh_ttl
}

func (cfg *Config) GetJWTSecret() string {
	return cfg.jwt_secret
}

func (cfg *Config) GetJWTttl() string {
	return cfg.jwt_ttl
}

func (cfg *Config) GetDBHost() string {
	return cfg.db_host
}

func (cfg *Config) GetDBHostPort() string {
	return cfg.db_port
}

func (cfg *Config) GetDBName() string {
	return cfg.db_name
}

func (cfg *Config) GetDBUsername() string {
	return cfg.db_user
}

func (cfg *Config) GetDBPassword() string {
	return cfg.db_password
}

func (cfg *Config) GetGRPCHost() string {
	return cfg.grpc_host
}

func (cfg *Config) GetGRPCPort() string {
	return cfg.grpc_port
}

func (cfg *Config) GetRESTHost() string {
	return cfg.rest_host
}

func (cfg *Config) GetRESTPort() string {
	return cfg.rest_port
}

func (cfg *Config) GetMainDir() string {
	return cfg.main_dir
}
