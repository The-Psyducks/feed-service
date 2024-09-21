package src

import (
	"log"
	"github.com/spf13/viper"

)

type Config struct {
	Enviroment string
	Port       string
	Host	   string
	Gin_Mode   string
	Mongo_URI  string
}

func ConfigEnv() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file, error: ", err)
	}

	viper.SetDefault("ENVIROMENT", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")

	config := &Config{
		Enviroment: viper.GetString("ENVIROMENT"),
		Port:       viper.GetString("PORT"),
		Host:       viper.GetString("HOST"),
		Gin_Mode:   viper.GetString("GIN_MODE"),
		Mongo_URI:  viper.GetString("MONGO_URI"),
	}

	return config
}