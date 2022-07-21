package auth

import "github.com/AscaroLabs/go-news/internal/config"

type RegistgationService struct {
}

func NewRegistgationService(cfg *config.Config) (*RegistgationService, error) {
	return &RegistgationService{}, nil
}
