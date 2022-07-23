package auth

import (
	"fmt"
	"log"

	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/AscaroLabs/go-news/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// регистрируем нового пользователя
func RegisterUser(cfg *config.Config, tm *TokenManager, userDTO *storage.UserDTO) (*Tokens, error) {

	hashedPassword, err := hashPassword(userDTO.Password)
	if err != nil {
		return nil, err
	}

	log.Printf("Let's make a txn!\n")

	err = storage.MakeRegistrationTxn(cfg, storage.UserDTO{
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: hashedPassword,
		Role:     userDTO.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("error while making txn: %v", err)
	}

	log.Printf("txn commited!\n")

	log.Printf("let's generate some tokens...\n")

	tokens, err := tm.GenerateTokensFromUserDTO(userDTO)
	if err != nil {
		return nil, fmt.Errorf("error while creating tokens: %v", err)
	}

	log.Printf("nice, tokens %s:%s", tokens.AccessToken, tokens.RefreshToken)
	return tokens, nil
}

func hashPassword(pas string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(pas), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), nil
}

func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}
