package db_test

// import (
// 	"log"
// 	"os"
// 	"testing"

// 	"region-llc-todo/pkg/config"
// 	"region-llc-todo/pkg/db"

// 	"github.com/spf13/viper"
// )

// var TestStorage db.Storage

// func TestMain(m *testing.M) {
// 	config, err := loadTestConfig()
// 	if err != nil {
// 		log.Fatal("cannot load config:", err)
// 	}

// 	TestStorage = db.Init(config.DBUrl)
// 	os.Exit(m.Run())
// }

// func loadTestConfig() (config config.Config, err error) {
// 	viper.SetDefault("PORT", ":50052")
// 	viper.SetDefault("DB_URL", "postgres://user:secret@localhost:5432/url_redirector")
// 	viper.SetDefault("JWT_SECRET_KEY", "not-secret-key")
// 	viper.SetDefault("ISSUER", "URL-svc")
// 	viper.SetDefault("EXPIRATION_HOURS", 1)

// 	viper.AutomaticEnv()
// 	if err = viper.Unmarshal(&config); err != nil {
// 		return
// 	}

// 	return
// }
