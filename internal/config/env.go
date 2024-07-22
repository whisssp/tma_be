package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

// Sup: supabase
type Config struct {
	Port             string
	Host             string
	DB               string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBAddress        string
	DBName           string
	RedisPort        string
	RedisPassword    string
	RedisDB          int
	SupUrl           string
	SupKey           string
	SupStorageRawUrl string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		Host:             getEnv("HOST", "http://localhost"),
		Port:             getEnv("PORT", "8080"),
		DB:               getEnv("DB", "postgres"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASS", "123456"),
		DBAddress:        fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "5432")),
		DBName:           getEnv("DB_NAME", "tma_db"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          int(getEnvAsInt("REDIS_DB", 0)),
		SupStorageRawUrl: getEnv("SUPABASE_STORAGE_RAW_URL", ""),
		SupUrl:           getEnv("SUPABASE_URL", ""),
		SupKey:           getEnv("SUPABASE_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}