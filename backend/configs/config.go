package config

import (
	"os"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
	JWTSecret  string
	UploadPath string

	AlipayAppID      string
	AlipayPrivateKey string
	AlipayPublicKey  string
	AlipayNotifyURL  string

	WxpayAppID     string
	WxpayMchID     string
	WxpayAPIKey    string
	WxpayCertPath  string
	WxpayKeyPath   string
	WxpayNotifyURL string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "force_learning"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		UploadPath: getEnv("UPLOAD_PATH", "./uploads"),

		AlipayAppID:      getEnv("ALIPAY_APP_ID", ""),
		AlipayPrivateKey: getEnv("ALIPAY_PRIVATE_KEY", ""),
		AlipayPublicKey:  getEnv("ALIPAY_PUBLIC_KEY", ""),
		AlipayNotifyURL:  getEnv("ALIPAY_NOTIFY_URL", ""),

		WxpayAppID:     getEnv("WXPAY_APP_ID", ""),
		WxpayMchID:     getEnv("WXPAY_MCH_ID", ""),
		WxpayAPIKey:    getEnv("WXPAY_API_KEY", ""),
		WxpayCertPath:  getEnv("WXPAY_CERT_PATH", ""),
		WxpayKeyPath:   getEnv("WXPAY_KEY_PATH", ""),
		WxpayNotifyURL: getEnv("WXPAY_NOTIFY_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
