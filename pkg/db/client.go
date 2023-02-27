package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
)

type pgConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func InitPGConfig() (*pgConfig, error) {
	var cfg pgConfig
	cfg.DBName = viper.GetString("db.name")
	cfg.Host = viper.GetString("db.host")
	cfg.Port = viper.GetString("db.port")
	cfg.Username = viper.GetString("db.user")
	cfg.Password = viper.GetString("db.password")

	if cfg.Host == "" || cfg.Port == "" ||
		cfg.Password == "" || cfg.Username == "" || cfg.DBName == "" {
		return nil, errors.New(fmt.Sprintf("error: value from config is null: %+v ", cfg))

	}
	return &cfg, nil
}

func GetPGClient(ctx context.Context, cfg *pgConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Println("error: connecting to DB")
		return nil, err
	}
	return pool, nil
}
