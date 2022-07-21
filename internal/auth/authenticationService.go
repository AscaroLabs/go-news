package auth

import "github.com/AscaroLabs/go-news/internal/config"

type AuthenticationService struct {
}

func NewAuthenticationService(cfg *config.Config) (*AuthenticationService, error) {
	return &AuthenticationService{}, nil
}
