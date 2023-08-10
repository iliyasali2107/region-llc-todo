package db_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"region-llc-todo-service/pkg/config"
	"region-llc-todo-service/pkg/db"

	"github.com/spf13/viper"
)

var TestStorage db.Storage

func TestMain(m *testing.M) {
	config, err := loadTestConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	TestStorage = db.Init(config)
	os.Exit(m.Run())
}

func loadTestConfig() (config config.Config, err error) {
	viper.SetDefault("DB_URL", "mongodb://localhost:27017/todo_list")
	viper.SetDefault("DB_NAME", "todo_list")
	viper.SetDefault("TODO_COLLECTION", "todo")

	viper.AutomaticEnv()
	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	fmt.Println(config)
	return
}
