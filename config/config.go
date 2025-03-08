package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server *Server
	Db     *Db
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
	MaxOpenConns int           `env:"DB_MAX_OPEN_CONNS" env-default:"30"`
	MaxIdleConns int           `env:"DB_MAX_IDLE_CONNS" env-default:"30"`
	MaxIdleTime  time.Duration `env:"DB_MAX_IDLE_TIME" env-default:"15m"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config/config.yaml", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
	//REDO: should we update?

}
