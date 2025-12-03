package main

import (
	"fmt"
	"log"

	"github.com/mariclezhang/vps_backend/internal/model"
	"github.com/mariclezhang/vps_backend/internal/util"
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
		MaxOpenConns: 100,
		MaxIdleConns: 10,
	}

	if err := db.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Starting to seed database...")

	// 清空现有数据（仅用于开发环境）
	if err := cleanDatabase(); err != nil {
		log.Fatalf("Failed to clean database: %v", err)
	}

	// 创建测试用户
	if err := seedUsers(); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	// 创建订阅套餐
	if err := seedSubscriptionPlans(); err != nil {
		log.Fatalf("Failed to seed subscription plans: %v", err)
	}

	// 创建节点
	if err := seedNodes(); err != nil {
		log.Fatalf("Failed to seed nodes: %v", err)
	}

	// 创建公告
	if err := seedAnnouncements(); err != nil {
		log.Fatalf("Failed to seed announcements: %v", err)
	}

	log.Println("Database seeded successfully!")
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func cleanDatabase() error {
	// 注意：这会删除所有数据！仅用于开发环境
	tables := []string{
		"user_node_access",
		"traffic_logs",
		"orders",
		"subscriptions",
		"subscription_plans",
		"nodes",
		"announcements",
		"password_resets",
		"users",
	}

	for _, table := range tables {
		if err := db.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			log.Printf("Warning: Failed to truncate %s: %v", table, err)
		}
	}

	return nil
}

func seedUsers() error {
	log.Println("Seeding users...")

	// 创建管理员用户
	adminPassword, _ := util.HashPassword("admin123456")
	admin := model.User{
		Email:        "admin@example.com",
		Username:     "admin",
		PasswordHash: adminPassword,
		Balance:      10000.00,
		Status:       "active",
	}
	if err := db.DB.Create(&admin).Error; err != nil {
		return err
	}

	// 创建测试用户
	testPassword, _ := util.HashPassword("123456")
	testUsers := []model.User{
		{
			Email:        "demo@example.com",
			Username:     "demo",
			PasswordHash: testPassword,
			Balance:      1000.00,
			Status:       "active",
		},
		{
			Email:        "user1@example.com",
			Username:     "user1",
			PasswordHash: testPassword,
			Balance:      500.00,
			Status:       "active",
		},
		{
			Email:        "user2@example.com",
			Username:     "user2",
			PasswordHash: testPassword,
			Balance:      200.00,
			Status:       "active",
		},
	}

	for _, user := range testUsers {
		if err := db.DB.Create(&user).Error; err != nil {
			return err
		}
	}

	log.Printf("Created %d users", len(testUsers)+1)
	return nil
}

func seedSubscriptionPlans() error {
	log.Println("Seeding subscription plans...")

	plans := []model.SubscriptionPlan{
		{
			Name:         "基础套餐",
			Description:  "适合轻度使用者",
			Price:        29.90,
			TrafficLimit: 100 * 1024 * 1024 * 1024, // 100GB
			DurationDays: 30,
			Features:     model.StringArray{"100GB流量", "5个设备同时在线", "标准速度"},
			IsActive:     true,
			SortOrder:    1,
		},
		{
			Name:         "标准套餐",
			Description:  "最受欢迎的选择",
			Price:        49.90,
			TrafficLimit: 200 * 1024 * 1024 * 1024, // 200GB
			DurationDays: 30,
			Features:     model.StringArray{"200GB流量", "10个设备同时在线", "高速连接", "优先支持"},
			IsActive:     true,
			SortOrder:    2,
		},
		{
			Name:         "高级套餐",
			Description:  "专业用户首选",
			Price:        99.90,
			TrafficLimit: 500 * 1024 * 1024 * 1024, // 500GB
			DurationDays: 30,
			Features:     model.StringArray{"500GB流量", "无限设备", "极速连接", "专属客服", "高级节点"},
			IsActive:     true,
			SortOrder:    3,
		},
		{
			Name:         "旗舰套餐",
			Description:  "无限制体验",
			Price:        199.90,
			TrafficLimit: 1024 * 1024 * 1024 * 1024, // 1TB
			DurationDays: 30,
			Features:     model.StringArray{"1TB流量", "无限设备", "极速连接", "专属客服", "所有节点", "优先网络"},
			IsActive:     true,
			SortOrder:    4,
		},
	}

	for _, plan := range plans {
		if err := db.DB.Create(&plan).Error; err != nil {
			return err
		}
	}

	log.Printf("Created %d subscription plans", len(plans))
	return nil
}

func seedNodes() error {
	log.Println("Seeding nodes...")

	nodes := []model.Node{
		{
			Name:               "美国-洛杉矶-01",
			Location:           "US",
			Protocol:           "vmess",
			ServerAddress:      "us-la-01.example.com",
			ServerPort:         443,
			Status:             "online",
			Latency:            120,
			LoadPercentage:     35,
			Bandwidth:          "1Gbps",
			MaxConnections:     1000,
			CurrentConnections: 350,
			IsActive:           true,
		},
		{
			Name:               "美国-西雅图-01",
			Location:           "US",
			Protocol:           "vmess",
			ServerAddress:      "us-sea-01.example.com",
			ServerPort:         443,
			Status:             "online",
			Latency:            110,
			LoadPercentage:     28,
			Bandwidth:          "1Gbps",
			MaxConnections:     1000,
			CurrentConnections: 280,
			IsActive:           true,
		},
		{
			Name:               "日本-东京-01",
			Location:           "JP",
			Protocol:           "vmess",
			ServerAddress:      "jp-tko-01.example.com",
			ServerPort:         443,
			Status:             "online",
			Latency:            60,
			LoadPercentage:     45,
			Bandwidth:          "1Gbps",
			MaxConnections:     1000,
			CurrentConnections: 450,
			IsActive:           true,
		},
		{
			Name:               "日本-大阪-01",
			Location:           "JP",
			Protocol:           "vmess",
			ServerAddress:      "jp-osa-01.example.com",
			ServerPort:         443,
			Status:             "online",
			Latency:            65,
			LoadPercentage:     40,
			Bandwidth:          "1Gbps",
			MaxConnections:     1000,
			CurrentConnections: 400,
			IsActive:           true,
		},
		{
			Name:               "新加坡-01",
			Location:           "SG",
			Protocol:           "vmess",
			ServerAddress:      "sg-01.example.com",
			ServerPort:         443,
			Status:             "online",
			Latency:            80,
			LoadPercentage:     50,
			Bandwidth:          "1Gbps",
			MaxConnections:     1000,
			CurrentConnections: 500,
			IsActive:           true,
		},
		{
			Name:               "香港-01",
			Location:           "HK",
			Protocol:           "vmess",
			ServerAddress:      "hk-01.example.com",
			ServerPort:         443,
			Status:             "online",
			Latency:            30,
			LoadPercentage:     60,
			Bandwidth:          "1Gbps",
			MaxConnections:     1000,
			CurrentConnections: 600,
			IsActive:           true,
		},
		{
			Name:               "德国-法兰克福-01",
			Location:           "DE",
			Protocol:           "vmess",
			ServerAddress:      "de-fra-01.example.com",
			ServerPort:         443,
			Status:             "online",
			Latency:            200,
			LoadPercentage:     30,
			Bandwidth:          "1Gbps",
			MaxConnections:     1000,
			CurrentConnections: 300,
			IsActive:           true,
		},
		{
			Name:               "英国-伦敦-01",
			Location:           "UK",
			Protocol:           "vmess",
			ServerAddress:      "uk-lon-01.example.com",
			ServerPort:         443,
			Status:             "online",
			Latency:            210,
			LoadPercentage:     25,
			Bandwidth:          "1Gbps",
			MaxConnections:     1000,
			CurrentConnections: 250,
			IsActive:           true,
		},
	}

	for _, node := range nodes {
		if err := db.DB.Create(&node).Error; err != nil {
			return err
		}
	}

	log.Printf("Created %d nodes", len(nodes))
	return nil
}

func seedAnnouncements() error {
	log.Println("Seeding announcements...")

	announcements := []model.Announcement{
		{
			Title:    "欢迎使用 VPS 管理平台",
			Content:  "感谢您选择我们的服务！如有任何问题，请随时联系客服。",
			Type:     "success",
			IsActive: true,
		},
		{
			Title:    "系统维护通知",
			Content:  "我们将于本周六凌晨2:00-4:00进行系统维护，期间服务可能短暂中断。",
			Type:     "warning",
			IsActive: true,
		},
		{
			Title:    "新增日本节点",
			Content:  "我们新增了大阪节点，欢迎体验更快的连接速度！",
			Type:     "info",
			Link:     "/nodes",
			IsActive: true,
		},
	}

	for _, announcement := range announcements {
		if err := db.DB.Create(&announcement).Error; err != nil {
			return err
		}
	}

	log.Printf("Created %d announcements", len(announcements))
	return nil
}
