package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Настройки пула (можно вынести в конфиг потом)
	config.MaxConns = 10                       // Максимум соединений
	config.MinConns = 2                        // Минимум (держать открытыми)
	config.MaxConnLifetime = 1 * time.Hour     // Время жизни соединения
	config.MaxConnIdleTime = 30 * time.Minute  // Время простоя

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}