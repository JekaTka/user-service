package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/JekaTka/user-service/internal/config"
)

var testStore Store

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("../../..")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	connPool, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
