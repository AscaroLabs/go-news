package main

import (
	"log"

	"github.com/AscaroLabs/go-news/internal/auth"
	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/AscaroLabs/go-news/internal/storage"
)

func main() {
	cfg := config.NewConfig()
	ok, err := storage.CheckHealth(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(ok)

	tm, err := auth.NewTokenManager(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// tokens, err := auth.RegisterUser(cfg, tm, &storage.UserDTO{
	// 	Name:     "Ilya",
	// 	Email:    "ad.ru",
	// 	Password: "qweasdzxc123123",
	// 	Role:     "admin",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("access %s\n refresh %s\n", tokens.AccessToken, tokens.RefreshToken)

	signin_tokens, err := auth.SignIn(cfg, tm, storage.SignInDTO{
		Email:    "ad@.ru",
		Password: "qweasdzxc123123",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("SignIn tokens:\nAccess: %s\nRefresh: %s", signin_tokens.AccessToken, signin_tokens.RefreshToken)
}
