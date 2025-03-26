package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var (
	connectErrorString        = "connect to database error: %v"
	migrateConnectErrorString = "connect for migrate database error: %v"
	migrateRunErrorString     = "run migrate error: %v"
)

type Config struct {
	Port uint16 `env:"POSTGRES_PORT" env-default:"5432"`
	Host string `env:"POSTGRES_HOST" env-default:"127.0.0.1"`

	Name     string `env:"POSTGRES_DB" env-default:"postgres"`
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`

	MaxConn int32 `env:"POSTGRES_MAX_CONN" env-default:"15"`
	MinConn int32 `env:"POSTGRES_MIN_CONN" env-default:"5"`
}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.MaxConn,
		cfg.MinConn)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf(connectErrorString, err)
		return nil, err
	}

	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		))
	if err != nil {
		log.Fatalf(migrateConnectErrorString, err)
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf(migrateRunErrorString, err)
		return nil, err
	}

	return pool, nil
}
