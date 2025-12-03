package main

import (
	"fmt"
	"log"

	"github.com/mariclezhang/vps_backend/internal/api/router"
	"github.com/mariclezhang/vps_backend/internal/util"
	"github.com/mariclezhang/vps_backend/pkg/cache"
	"github.com/mariclezhang/vps_backend/pkg/db"
	"github.com/spf13/viper"
)

func main() {
	// 加载配置
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
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

	if err := db.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")

	// 初始化Redis
	redisConfig := cache.Config{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetInt("redis.port"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}

	if err := cache.InitRedis(redisConfig); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}

	// 初始化JWT
	util.InitJWT(viper.GetString("jwt.secret"))

	// 设置路由
	frontendURL := viper.GetString("server.frontend_url")
	r := router.SetupRouter(frontendURL)

	// 启动服务器
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// loadConfig 加载配置文件
func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 读取环境变量
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.frontend_url", "http://localhost:8000")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.expire_hours", 24)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults")
			return nil
		}
		return err
	}

	log.Println("Config file loaded successfully")
	return nil
}
