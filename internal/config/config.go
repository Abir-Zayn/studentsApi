package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// @struct tags:yaml
// env-default:"production"
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" `
	StoragePath string `yaml:"storage_path"`
	HTTPServer 	`yaml:"http_server"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	//Stat method 
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist:%s", configPath)
	}

	var cfg Config 
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Failed to read config file:%s, error:%s", configPath, err.Error())
	}
	
	return &cfg


}