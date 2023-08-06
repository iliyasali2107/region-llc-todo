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
	ServicePort        string `mapstructure:"SERVICE_PORT"`
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

// func LoadConfig() (config Config, err error) {
// viper.SetDefault("SERVICE_PORT", "todo_service:5000")
// 	viper.SetDefault("PORT", ":4000")
// 	viper.SetDefault("DB_URL", "mongodb://mongo:27017/mymongo")
// 	viper.SetDefault("JWT_SECRET_KEY", "not-secret-key")
// 	viper.SetDefault("ISSUER", "")
// 	viper.SetDefault("EXPIRATION_HOURS", 1)
// 	viper.SetDefault("CLIENT_PORT", ":4000")
// 	viper.SetDefault("DB_NAME", "todo_list")
// 	viper.SetDefault("TODO_COLLECTION", "todo")

// 	viper.AutomaticEnv()
// 	if err = viper.Unmarshal(&config); err != nil {
// 		return
// 	}

// 	return
// }

func LoadConfig() (config Config, err error) {
	viper.SetDefault("SERVICE_PORT", ":5000")
	viper.SetDefault("DB_URL", "mongodb://localhost:27017/mymongo")
	viper.SetDefault("JWT_SECRET_KEY", "not-secret-key")
	viper.SetDefault("ISSUER", "")
	viper.SetDefault("EXPIRATION_HOURS", 1)
	viper.SetDefault("CLIENT_PORT", ":4000")
	viper.SetDefault("DB_NAME", "mymongo")
	viper.SetDefault("TODO_COLLECTION", "todo")

	viper.AutomaticEnv()
	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return
}
