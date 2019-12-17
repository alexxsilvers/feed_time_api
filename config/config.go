package config

import (
	"context"
	"sync"

	"github.com/caarlos0/env"

	"github.com/alexxsilvers/feed_time_api/core/logger"
)

var mu sync.Mutex
var cfg config

type config struct {
	GracefulDelayInSeconds     int    `env:"GRACEFUL_DELAY_IN_SECONDS"`
	GracefulTimeoutInSeconds   int    `env:"GRACEFUL_TIMEOUT_IN_SECONDS"`
	PublicPort                 uint   `env:"PUBLIC_PORT"`
	CarsApiEndpoint            string `env:"CARS_API_ENDPOINT"`
	CarsApiTimeoutInSeconds    int    `env:"CARS_API_TIMEOUT_IN_SECONDS"`
	PredictApiEndpoint         string `env:"PREDICT_API_ENDPOINT"`
	PredictApiTimeoutInSeconds int    `env:"PREDICT_API_TIMEOUT_IN_SECONDS"`
}

func ParseConfig() {
	conf := config{}
	err := env.Parse(&conf)
	if err != nil {
		logger.Fatal(context.Background(), err)
	}
	cfg = conf
}

func GetGracefulDelayInSeconds() int {
	mu.Lock()
	val := cfg.GracefulTimeoutInSeconds
	mu.Unlock()
	return val
}

func GetGracefulTimeoutInSeconds() int {
	mu.Lock()
	val := cfg.GracefulTimeoutInSeconds
	mu.Unlock()
	return val
}

func GetPublicPort() uint {
	mu.Lock()
	val := cfg.PublicPort
	mu.Unlock()
	return val
}

func GetCarsApiEndpoint() string {
	mu.Lock()
	val := cfg.CarsApiEndpoint
	mu.Unlock()
	return val
}

func GetCarsApiTimeoutInSeconds() int {
	mu.Lock()
	val := cfg.CarsApiTimeoutInSeconds
	mu.Unlock()
	return val
}

func GetPredictApiEndpoint() string {
	mu.Lock()
	val := cfg.PredictApiEndpoint
	mu.Unlock()
	return val
}

func GetPredictApiTimeoutInSeconds() int {
	mu.Lock()
	val := cfg.PredictApiTimeoutInSeconds
	mu.Unlock()
	return val
}
