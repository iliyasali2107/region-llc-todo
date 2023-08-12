package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	DBUrl              string `mapstructure:"DB_URL"`
	DbName             string `mapstructure:"DB_NAME"`
	TodoCollectionName string `mapstructure:"TODO_COLLECTION"`
}

// func LoadConfig() (config Config, err error) {
// 	viper.SetDefault("PORT", "todo_service:50051")
// 	viper.SetDefault("DB_URL", "mongodb://mongo:27017/todo_list")
// 	viper.SetDefault("DB_NAME", "todo_list")
// 	viper.SetDefault("TODO_COLLECTION", "todo")

// 	viper.AutomaticEnv()
// 	if err = viper.Unmarshal(&config); err != nil {
// 		return
// 	}

// 	return
// }

// for local
func LoadConfig() (config Config, err error) {
	viper.SetDefault("PORT", ":50051")
	viper.SetDefault("DB_URL", "mongodb://localhost:27017/todo_list")
	viper.SetDefault("DB_NAME", "todo_list")
	viper.SetDefault("TODO_COLLECTION", "todo")

	viper.AutomaticEnv()
	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return
}
