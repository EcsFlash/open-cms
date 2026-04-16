package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Database   Database   `yaml:"database"`
	JWT        JWT        `yaml:"jwt"`
	Admin      AdminData  `yaml:"admin"`
	MinIO      MinIO      `yaml:"minio"`
}

type HTTPServer struct {
	Address string `yaml:"address" env:"HTTP_ADDRESS" env-default:":8080"`
	SelfURL string `yaml:"self_url" env:"SELF_URL" env-default:"http://localhost:8080"`
}

type AdminData struct {
	Nickname string `yaml:"nickname"`
	Password string `yaml:"passwd"`
}
type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int32  `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"cms"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

type JWT struct {
	Secret string `yaml:"secret" env:"JWT_SECRET" env-default:"dev-secret-change-me"`
}

type MinIO struct {
	Endpoint     string `yaml:"endpoint" env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
	AccessKey    string `yaml:"access_key" env:"MINIO_ACCESS_KEY" env-default:"minioadmin"`
	SecretKey    string `yaml:"secret_key" env:"MINIO_SECRET_KEY" env-default:"minioadmin"`
	UseSSL       bool   `yaml:"use_ssl" env:"MINIO_USE_SSL" env-default:"false"`
	ImagesBucket string `yaml:"images_bucket" env:"MINIO_IMAGES_BUCKET" env-default:"cms-images"`
	VideosBucket string `yaml:"videos_bucket" env:"MINIO_VIDEOS_BUCKET" env-default:"cms-videos"`
}

func (d Database) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		d.Host, d.User, d.Password, d.DBName, d.Port, d.SSLMode,
	)
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	// Allow env vars to override YAML values (useful for docker-compose).
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read env: %s", err)
	}
	return &cfg
}
