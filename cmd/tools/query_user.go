package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mariclezhang/vps_backend/internal/model"
	"github.com/mariclezhang/vps_backend/pkg/db"
	"github.com/spf13/viper"
)

func main() {
	emailPtr := flag.String("email", "", "Email of the user to query")
	flag.Parse()

	if *emailPtr == "" {
		// Try to read from positional argument if flag is not set
		if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
			*emailPtr = os.Args[1]
		} else {
			fmt.Println("Usage: go run cmd/tools/query_user.go -email <email>")
			os.Exit(1)
		}
	}

	// Load Config
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init DB
	dbConfig := db.Config{
		Host:         viper.GetString("database.host"),
		Port:         viper.GetInt("database.port"),
		User:         viper.GetString("database.user"),
		Password:     viper.GetString("database.password"),
		DBName:       viper.GetString("database.dbname"),
		SSLMode:      viper.GetString("database.sslmode"),
		MaxOpenConns: viper.GetInt("database.max_open_conns"),
		MaxIdleConns: viper.GetInt("database.max_idle_conns"),
	}

	// Suppress GORM default logger for cleaner output, or keep it for debugging
	// main.go uses Info mode. We might want Silent or Error for this tool to just show JSON.
	// But db.InitDB sets it to Info. Let's stick to default.
	if err := db.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	var user model.User
	// Query
	if err := db.DB.Where("email = ?", *emailPtr).First(&user).Error; err != nil {
		log.Fatalf("Error finding user with email %s: %v", *emailPtr, err)
	}

	// Print JSON
	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling user: %v", err)
	}
	fmt.Println(string(b))
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// Look for config in project root/config
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults (Partial set, enough for DB if config file missing)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults")
			return nil
		}
		return err
	}
	return nil
}
