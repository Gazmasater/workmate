package config

import (
	"os"
	"strconv"
	"time"

	"github.com/gaz358/myprog/workmate/pkg/logger"
	"github.com/joho/godotenv"
)

const (
	defaultTaskDuration    = 60 * time.Second
	defaultShutdownTimeout = 5 * time.Second
)

type Config struct {
	Port            string
	LogLevel        string
	TaskDuration    time.Duration
	ShutdownTimeout time.Duration
}

func Load() *Config {
	log := logger.Global().Named("config")
	log.Info("loading config...")

	if err := godotenv.Load(); err != nil {
		log.Warn("[config] .env не найден, используются переменные окружения по умолчанию")
	} else {
		log.Info("[config] .env успешно загружен")
	}

	cfg := &Config{
		Port:            getEnv(log, "PORT", "8080"),
		LogLevel:        getEnv(log, "LOG_LEVEL", "info"),
		TaskDuration:    getEnvAsDuration(log, "TASK_DURATION", defaultTaskDuration),
		ShutdownTimeout: getEnvAsDuration(log, "SHUTDOWN_TIMEOUT", defaultShutdownTimeout),
	}

	log.Infow("[config] итоговая конфигурация",
		"PORT", cfg.Port,
		"LOG_LEVEL", cfg.LogLevel,
		"TASK_DURATION", cfg.TaskDuration,
		"SHUTDOWN_TIMEOUT", cfg.ShutdownTimeout,
	)

	return cfg
}

func getEnv(log logger.TypeOfLogger, key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Infow("[config] переменная не установлена, используется по умолчанию", "key", key, "default", def)
		return def
	}
	log.Debugw("[config] переменная окружения", "key", key, "value", val)
	return val
}

func getEnvAsDuration(log logger.TypeOfLogger, key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		log.Infow("[config] переменная не установлена, используется значение по умолчанию", "key", key, "default", def)
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Warnw("[config] неверное значение, используется по умолчанию", "key", key, "value", val, "err", err, "default", def)
		return def
	}
	log.Debugw("[config] переменная окружения (duration)", "key", key, "seconds", i)
	return time.Duration(i) * time.Second
}
