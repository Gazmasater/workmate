package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	LogLevel        string
	TaskDuration    time.Duration
	ShutdownTimeout time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env не найден, используются переменные окружения по умолчанию")
	}

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		TaskDuration:    getEnvAsDuration("TASK_DURATION", 60*time.Second),
		ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
	}

	log.Printf("[config] PORT=%s", cfg.Port)
	log.Printf("[config] LOG_LEVEL=%s", cfg.LogLevel)
	log.Printf("[config] TASK_DURATION=%s", cfg.TaskDuration)
	log.Printf("[config] SHUTDOWN_TIMEOUT=%s", cfg.ShutdownTimeout)

	return cfg
}

func getEnv(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getEnvAsDuration(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("[config] неверное значение %s=%q: %v, используется по умолчанию: %v", key, val, err, def)
		return def
	}
	return time.Duration(i) * time.Second
}
