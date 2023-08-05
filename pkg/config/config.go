package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	DBUrl              string `mapstructure:"DB_URL"`
	JWTSecretKey       string `mapstructure:"JWT_SECRET_KEY"`
	Issuer             string `mapstructure:"ISSUER"`
	ExpirationHours    int    `mapstructure:"EXPIRATION_HOURS"`
	ClientPort         string `mapstructure:"CLIENT_PORT"`
	DbName             string `mapstructure:"DB_NAME"`
	TodoCollectionName string `mapstructure:"TODO_COLLECTION"`
}

// func LoadConfig() (config Config, err error) {
// 	viper.SetDefault("PORT", "url_service:50051")
// 	viper.SetDefault("DB_URL", "postgres://user:secret@postgres:5432/url_redirector")
// 	viper.SetDefault("JWT_SECRET_KEY", "not-secret-key")
// 	viper.SetDefault("ISSUER", "URL-svc")
// 	viper.SetDefault("EXPIRATION_HOURS", 1)

// 	viper.AutomaticEnv()
// 	if err = viper.Unmarshal(&config); err != nil {
// 		return
// 	}

// 	return
// }

func LoadConfig() (config Config, err error) {
	viper.AutomaticEnv()                          // Read environment variables
	viper.SetConfigFile("./pkg/config/envs/.env") // Optionally path to look for config file
	viper.ReadInConfig()                          // Read the config file

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}
