package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/AscaroLabs/go-news/internal/storage"
	jwt "github.com/golang-jwt/jwt/v4"
)

type TokenManager struct {
	secret []byte
	method jwt.SigningMethod
	ttl    time.Duration
	cfg    *config.Config
}

// структура представляет собой пару токенов, готовых к отправке
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// функция генерирует пару токенов
func (tm *TokenManager) GenerateTokens(tknDTO *storage.TokenDTO) (*Tokens, error) {

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  tknDTO.UserId,
		"role": tknDTO.Role,
		"exp":  fmt.Sprintf("%d", time.Now().Add(tm.ttl).Unix()),
	})

	token, err := tkn.SignedString(tm.secret)
	if err != nil {
		return nil, err
	}
	return &Tokens{
		AccessToken:  token,
		RefreshToken: "refreshTokenExample",
	}, nil
}

// функция генерирует пару токенов по информации и пользователе
func (tm *TokenManager) GenerateTokensFromUserDTO(userDTO *storage.UserDTO) (*Tokens, error) {

	log.Printf("ok, we have UserDTO...\n")

	tokenDTO, err := storage.GetTokenDTOFromUserDTO(tm.cfg, userDTO)
	if err != nil {
		return nil, err
	}

	log.Printf("now we have TokenDTO( userId: %s, role: %s)\n", tokenDTO.UserId, tokenDTO.Role)

	log.Printf("let's generate tokens from this DTO!\n")

	tokens, err := tm.GenerateTokens(tokenDTO)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// функция парсит токен
func (tm *TokenManager) ParseToken(tkn string) (storage.TokenDTO, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("tnexpected signing method: %v", token.Header["alg"])
		}
		if claims, ok := token.Claims.(jwt.MapClaims); !ok || claims["exp"].(int64) > time.Now().Unix() {
			return nil, fmt.Errorf("token expired")
		}
		return tm.secret, nil
	})
	if err != nil {
		return storage.TokenDTO{}, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	return storage.TokenDTO{
		UserId: claims["sub"].(string),
		Role:   claims["role"].(string),
	}, nil
}

func NewTokenManager(cfg *config.Config) (*TokenManager, error) {
	jwt_ttl, err := time.ParseDuration(cfg.GetJWTttl())
	if err != nil {
		return nil, err
	}
	return &TokenManager{
		secret: []byte(cfg.GetJWTSecret()),
		method: jwt.SigningMethodHS256,
		ttl:    jwt_ttl,
		cfg:    cfg,
	}, nil
}
