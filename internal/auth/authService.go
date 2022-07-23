package auth

import (
	"log"

	"github.com/AscaroLabs/go-news/internal/config"
)

type AuthSerivce struct {
	tokenManager *TokenManager
}

func NewAuthService(cfg *config.Config) (*AuthSerivce, error) {
	tokenManager, err := NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("Can't create AuthService: %v", err)
	}
	return &AuthSerivce{
		tokenManager: tokenManager,
	}, nil
}
