package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server *Server
	Jwt    *Jwt
	Db     *Db
}
type Jwt struct {
	Signkey  string        `yaml:"signkey" env:"JWT_SIGNKEY" env-default:"secret"`
	TokenTTL time.Duration `yaml:"tokenTTL" env:"JWT_TOKENTTL" env-default:"5m"`
}
type Server struct {
	Port string `yaml:"port" env:"DB_PORT" env-default:":8080"`
	Host string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
}

type Db struct {
	Port         string        `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Host         string        `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	DbName       string        `yaml:"dbname" env:"DB_NAME" env-default:"shop"`
	User         string        `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password     string        `yaml:"password" env:"DB_PASSWORD" env-default:"password"`
	SSLMode      string        `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
	MaxOpenConns int           `yaml:"max_open_conn" env:"DB_MAX_OPEN_CONNS" env-default:"30"`
	MaxIdleConns int           `yaml:"max_idle_conn" env:"DB_MAX_IDLE_CONNS" env-default:"30"`
	MaxIdleTime  time.Duration `yaml:"max_idle_time" env:"DB_MAX_IDLE_TIME" env-default:"15m"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./internal/app/config/config.yaml", cfg); err != nil {
		return nil, err
	}
	dbHost := os.Getenv("DATABASE_HOST")
	if dbHost != "" {
		cfg.Db.Host = dbHost
	}

	return cfg, nil

}
