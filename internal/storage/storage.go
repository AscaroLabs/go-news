package storage

import (
	"context"
	"fmt"

	"github.com/AscaroLabs/go-news/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func getUrl(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.GetDBUsername(),
		cfg.GetDBPassword(),
		cfg.GetDBHost(),
		cfg.GetDBHostPort(),
		cfg.GetDBName(),
	)
}

func getPool(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), getUrl(cfg))
	return pool, err
}

// проверяем жизнеспособнось БД
func CheckHealth(cfg *config.Config) (bool, error) {
	pool, err := getPool(cfg)

	if err != nil {
		return false, err
	}
	defer pool.Close()

	var info string
	err = pool.QueryRow(context.Background(), "select 'ok'").Scan(&info)
	if err != nil {
		return false, err
	}
	return true, nil
}
