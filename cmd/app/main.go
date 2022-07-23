// создание объекта структуры конфига и запуск приложения
package main

import (
	"log"

	"github.com/AscaroLabs/go-news/internal/app"
	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/xlab/closer"
)

func main() {

	// функция, которая вызовется при получении приложением сигнала об остановке
	closer.Bind(func() {
		log.Print("Stop running...")
	})
	cfg := config.NewConfig()
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Can't create new app.App object: %v", err)
	}
	go application.Run()
	// ждем в основной горутине сигнал
	closer.Hold()
}
