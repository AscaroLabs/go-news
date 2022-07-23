package auth

import (
	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/AscaroLabs/go-news/internal/storage"
)

// должны взять мыло и пароль и выплюнуть челу токены

func SignIn(cfg *config.Config, tm *TokenManager, signInDTO storage.SignInDTO) (*Tokens, error) {
	userDTO, err := storage.GetUserDTObyEmail(cfg, signInDTO.Email)
	if err != nil {
		return nil, err
	}
	ok := doPasswordsMatch(userDTO.Password, signInDTO.Password)
	if !ok {
		return nil, err
	}
	tokens, err := tm.GenerateTokensFromUserDTO(userDTO)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
