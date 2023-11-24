package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

//описываем структуру как в config файле yaml
type Config struct {
	Env string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	//можно будет запускать go run main.go --config=./config/local.yaml не задавая env

	if res == "" {
		godotenv.Load()
		res = os.Getenv("CONFIG_PATH")
	}

	return  res
}

func MustLoadPath(configPath string) *Config {
	//проверяем что по пути что-то есть
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	//парсим файл и записываем его в конфиг
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}