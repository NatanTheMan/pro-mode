// Package config load the enviorment configs
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Token,
	DatabaseURL string
	DevID,
	ChatID,
	AdminID int
}

func Load() (*Config, error) {
	var errs []error

	devID, err := intEnv("DEV_ID")
	if err != nil {
		errs = append(errs, err)
	}
	adminID, err := intEnv("ADMIN_ID")
	if err != nil {
		errs = append(errs, err)
	}
	chatID, err := intEnv("CHAT_ID")
	if err != nil {
		errs = append(errs, err)
	}

	cfg := &Config{
		DevID:       devID,
		Token:       os.Getenv("TOKEN"),
		ChatID:      chatID,
		AdminID:     adminID,
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.Token == "" {
		errs = append(errs, errors.New("token is required"))
	}
	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("database connection string is required"))
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return cfg, nil
}

func intEnv(key string) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, fmt.Errorf("empty env: %s", key)
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid int: %s", key, value)
	}

	return i, nil
}
