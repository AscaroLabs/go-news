package auth

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/AscaroLabs/go-news/internal/storage"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenManager struct {
	secret []byte
	method jwt.SigningMethod
	ttl    time.Duration
	cfg    *config.Config
}

type key int

const tokenManagerKey key = 0

func NewContext(ctx context.Context, tm *TokenManager) context.Context {
	return context.WithValue(ctx, tokenManagerKey, tm)
}

func FromContext(ctx context.Context) (*TokenManager, bool) {
	tm, ok := ctx.Value(tokenManagerKey).(*TokenManager)
	return tm, ok
}

// структура представляет собой пару токенов, готовых к отправке
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// функция генерирует пару токенов
func (tm *TokenManager) GenerateTokens(tknDTO *storage.TokenDTO) (*Tokens, error) {

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      tknDTO.UserId,
		"name":     tknDTO.Name,
		"role":     tknDTO.Role,
		"expireat": fmt.Sprintf("%d", time.Now().Add(tm.ttl).Unix()),
	})

	token, err := tkn.SignedString(tm.secret)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()

	err = storage.MakeNewSession(tm.cfg, tknDTO.UserId, refreshToken)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  token,
		RefreshToken: refreshToken,
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
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("bad token")
		}

		log.Printf("Hey yo, claims: %v", claims)
		exp, err := strconv.Atoi(claims["expireat"].(string))
		log.Print("Yo hay")
		log.Printf("%v", time.Unix(int64(exp), 0))

		if err != nil {
			return nil, err
		}

		log.Printf("%v", int64(exp) < time.Now().Unix())
		if int64(exp) < time.Now().Unix() {

			log.Printf("%d < %d", int64(exp), time.Now().Unix())
			return nil, fmt.Errorf("expired token")
		}

		return tm.secret, nil
	})
	if err != nil {
		log.Printf("ParseToken error: %s", err.Error())
		return storage.TokenDTO{}, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	return storage.TokenDTO{
		UserId: claims["sub"].(string),
		Name:   claims["name"].(string),
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
