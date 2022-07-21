package auth

import "github.com/AscaroLabs/go-news/internal/config"

type AuthorizationService struct {
}

func NewAuthorizationService(cfg *config.Config) (*AuthorizationService, error) {
	return &AuthorizationService{}, nil
}
