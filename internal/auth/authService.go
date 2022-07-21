package auth

import (
	"log"

	"github.com/AscaroLabs/go-news/internal/config"
)

type AuthSerivce struct {
	tokenManager          *TokenManager
	registgationService   *RegistgationService
	authenticationService *AuthenticationService
	authorizationService  *AuthorizationService
}

func NewAuthService(cfg *config.Config) (*AuthSerivce, error) {
	tokenManager, err := NewTokenManager(cfg)
	if err != nil {
		log.Fatalf("Can't create AuthService: %v", err)
	}
	registgationService, err := NewRegistgationService(cfg)
	if err != nil {
		log.Fatalf("Can't create AuthService: %v", err)
	}
	authenticationService, err := NewAuthenticationService(cfg)
	if err != nil {
		log.Fatalf("Can't create AuthService: %v", err)
	}
	authorizationService, err := NewAuthorizationService(cfg)
	if err != nil {
		log.Fatalf("Can't create AuthService: %v", err)
	}
	return &AuthSerivce{
		tokenManager:          tokenManager,
		registgationService:   registgationService,
		authenticationService: authenticationService,
		authorizationService:  authorizationService,
	}, nil
}
