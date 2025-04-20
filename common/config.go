package common

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"port"`
	DatabaseURL string `mapstructure:"database_url"`
}

func LoadConfig() *Config {
	var config *Config
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatalf("Fatal error Config file: %s \n", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {               // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	return config
}
