package auth

import "github.com/AscaroLabs/go-news/internal/config"

type TokenManager struct {
}

func NewTokenManager(cfg *config.Config) (*TokenManager, error) {
	return &TokenManager{}, nil
}
