package config

import "os"


// Config struct
type Config struct {
	Enviroment string
	Port       string
	Host       string
	Gin_Mode   string
	Mongo_URI  string
}

// ConfigEnv function
func ConfigEnv() *Config {
	config := &Config{
		Enviroment: getEnvOrDefault("ENVIROMENT", "development"),
		Port:      	getEnvOrDefault("PORT", "8080"),
		Host:       getEnvOrDefault("HOST", "0.0.0.0"),
		Gin_Mode:   getEnvOrDefault("GIN_MODE", "debug"),
		Mongo_URI:  getEnvOrDefault("MONGO_URI", "mongodb://mongo:27017"),
	}

	if os.Getenv("ENVIROMENT") == "HEROKU" {
		config.Mongo_URI = os.Getenv("DATABASE_URL")
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}