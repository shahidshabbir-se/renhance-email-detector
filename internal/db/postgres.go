package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shahidshabbir-se/renhance-email-detector/pkg/utils"
	"github.com/sirupsen/logrus"
)

var PGPool *pgxpool.Pool

func InitPostgres(ctx context.Context, log *logrus.Logger) *pgxpool.Pool {
	dsn := utils.GetEnv("DATABASE_URL", "postgres://harald_user:harald_pass@localhost:5432/emaildb")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.WithError(err).Fatal("PostgreSQL pool creation failed")
	}

	if err := pool.Ping(ctx); err != nil {
		log.WithError(err).Fatal("PostgreSQL ping failed")
	}

	log.Info("PostgreSQL connected")
	PGPool = pool
	return pool
}
