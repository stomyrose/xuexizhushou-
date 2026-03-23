package main

import (
	"fmt"
	"log"

	"force-learning/configs"
	"force-learning/internal/api/handler"
	"force-learning/internal/api/router"
	"force-learning/internal/model"
	"force-learning/internal/repository"
	"force-learning/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.SubscriptionPlan{},
		&model.Subscription{},
		&model.KnowledgeFile{},
		&model.LearningRecord{},
		&model.PayOrder{},
		&model.SyncRecord{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedSubscriptionPlans(db)

	_ = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	userRepo := repository.NewUserRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	learningRepo := repository.NewLearningRecordRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, userRepo)
	knowledgeService := service.NewKnowledgeService(knowledgeRepo, cfg.UploadPath)
	learningService := service.NewLearningService(learningRepo, userRepo)

	alipayConfig := &service.AlipayConfig{
		AppID:           cfg.AlipayAppID,
		PrivateKey:      cfg.AlipayPrivateKey,
		AlipayPublicKey: cfg.AlipayPublicKey,
		NotifyURL:       cfg.AlipayNotifyURL,
	}

	wxpayConfig := &service.WxpayConfig{
		AppID:     cfg.WxpayAppID,
		MchID:     cfg.WxpayMchID,
		APIKey:    cfg.WxpayAPIKey,
		CertPath:  cfg.WxpayCertPath,
		KeyPath:   cfg.WxpayKeyPath,
		NotifyURL: cfg.WxpayNotifyURL,
	}

	paymentService := service.NewPaymentService(subscriptionRepo, userRepo, alipayConfig, wxpayConfig)
	syncService := service.NewSyncService(learningRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, paymentService)
	knowledgeHandler := handler.NewKnowledgeHandler(knowledgeService)
	learningHandler := handler.NewLearningHandler(learningService, syncService)

	r := router.NewRouter(
		authHandler,
		subscriptionHandler,
		knowledgeHandler,
		learningHandler,
		cfg.JWTSecret,
	)

	engine := r.Setup()

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func seedSubscriptionPlans(db *gorm.DB) {
	var count int64
	db.Model(&model.SubscriptionPlan{}).Count(&count)
	if count > 0 {
		return
	}

	plans := []model.SubscriptionPlan{
		{Name: "月度订阅", DurationDays: 30, Price: 29.9, IsActive: true},
		{Name: "季度订阅", DurationDays: 90, Price: 79.9, IsActive: true},
		{Name: "年度订阅", DurationDays: 365, Price: 299.9, IsActive: true},
	}

	for _, plan := range plans {
		db.Create(&plan)
	}
}

var _ = gin.Default
