package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/db/sqlc"
)

type Store interface {
	sqlc.Querier
	WithTx(tx pgx.Tx) Store
}

type PGXStore struct {
	*sqlc.Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &PGXStore{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (s *PGXStore) WithTx(tx pgx.Tx) Store {
	return &PGXStore{
		db:      s.db,
		Queries: sqlc.New(tx),
	}
}
