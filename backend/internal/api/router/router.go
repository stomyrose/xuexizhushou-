package router

import (
	"force-learning/internal/api/handler"
	"force-learning/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine              *gin.Engine
	authHandler         *handler.AuthHandler
	subscriptionHandler *handler.SubscriptionHandler
	knowledgeHandler    *handler.KnowledgeHandler
	learningHandler     *handler.LearningHandler
	jwtSecret           string
}

func NewRouter(
	authHandler *handler.AuthHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	knowledgeHandler *handler.KnowledgeHandler,
	learningHandler *handler.LearningHandler,
	jwtSecret string,
) *Router {
	return &Router{
		authHandler:         authHandler,
		subscriptionHandler: subscriptionHandler,
		knowledgeHandler:    knowledgeHandler,
		learningHandler:     learningHandler,
		jwtSecret:           jwtSecret,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine = gin.Default()

	r.engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Admin-Token")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	v1 := r.engine.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.Refresh)
			auth.POST("/verify", middleware.Auth(r.jwtSecret), r.authHandler.Verify)
			auth.GET("/status", middleware.Auth(r.jwtSecret), r.authHandler.GetStatus)
		}

		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.GET("/plans", r.subscriptionHandler.GetPlans)
			subscriptions.POST("/purchase", middleware.Auth(r.jwtSecret), r.subscriptionHandler.Purchase)
			subscriptions.GET("/current", middleware.Auth(r.jwtSecret), r.subscriptionHandler.GetCurrent)
		}

		knowledge := v1.Group("/knowledge")
		{
			knowledge.GET("/files", middleware.Auth(r.jwtSecret), r.knowledgeHandler.ListFiles)
			knowledge.GET("/random", middleware.Auth(r.jwtSecret), r.knowledgeHandler.GetRandom)
			knowledge.GET("/download/:id", middleware.Auth(r.jwtSecret), r.knowledgeHandler.Download)
			knowledge.POST("/upload", middleware.Auth(r.jwtSecret), middleware.Admin(), r.knowledgeHandler.Upload)
			knowledge.DELETE("/files/:id", middleware.Auth(r.jwtSecret), middleware.Admin(), r.knowledgeHandler.Delete)
		}

		learning := v1.Group("/learning")
		{
			learning.POST("/records", middleware.Auth(r.jwtSecret), r.learningHandler.CreateRecord)
			learning.GET("/records", middleware.Auth(r.jwtSecret), r.learningHandler.GetRecords)
			learning.GET("/statistics", middleware.Auth(r.jwtSecret), r.learningHandler.GetStatistics)
			learning.POST("/records/batch", middleware.Auth(r.jwtSecret), r.learningHandler.BatchCreate)
		}
	}

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}
