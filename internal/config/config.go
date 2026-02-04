package config

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App    AppConfig
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

type AppConfig struct {
	Env string `env:"ENV" env-default:"local"`
}

type ServerConfig struct {
	Port            int           `env-required:"true" yaml:"port" env:"APP_PORT"`
	ReadTimeout     time.Duration `env-required:"true" yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `env-required:"true" yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `env-required:"true" yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
}

type DBConfig struct {
	Host     string `env-required:"true" env:"DB_HOST"`
	Port     string `env-required:"true" env:"DB_PORT"`
	User     string `env-required:"true" env:"DB_USER"`
	Password string `env-required:"true" env:"DB_PASSWORD"`
	Name     string `env-required:"true" env:"DB_NAME"`
	PoolSize int    `env-required:"true" yaml:"pool_size" env:"DB_POOL_SIZE"`
}

func (db *DBConfig) Url() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}

func New(configPath string) *Config {
	var config Config

	err := cleanenv.ReadConfig(path.Join("./", configPath), &config)
	if err != nil {
		panic(fmt.Sprintf("config:New:ReadConfig - %s", err.Error()))
	}

	return &config
}
