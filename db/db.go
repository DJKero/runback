package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

func New() {
	var err error
	DBPool, err = pgxpool.New(context.Background(), "postgresql://kero:kvAaVlB9a5O3LZ9xglHPf1NXtwTmcz@localhost:5432/runback?sslmode=disable")
	if err != nil {
		log.Fatalln("[ERROR] Error creating DB connection pool,", err)
	}
	_ = DBPool
}
