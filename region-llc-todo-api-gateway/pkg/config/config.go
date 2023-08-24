package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port            string `mapstructure:"PORT"`
	TodoServicePort string `mapstructure:"TODO_SERVICE_PORT"`
}

func LoadConfig() (config Config, err error) {
	viper.SetDefault("TODO_SERVICE_PORT", "todo_service:50051")
	viper.SetDefault("PORT", ":4000")

	viper.AutomaticEnv()
	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return
}

// LOCAL
// func LoadConfig() (config Config, err error) {
// 	viper.SetDefault("TODO_SERVICE_PORT", ":50051")
// 	viper.SetDefault("PORT", ":4000")

// 	viper.AutomaticEnv()
// 	if err = viper.Unmarshal(&config); err != nil {
// 		return
// 	}

// 	return
// }
