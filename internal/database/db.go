package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

package database

import (
"github.com/jackc/pgx/v4/pgxpool"
"golang.org/x/net/context"
)

type DataBase struct {
	DB *pgxpool.Pool
}

func CreateDb(connection string) (*DataBase, error) {

	pool, err := pgxpool.Connect(context.Background(), connection)
	if err != nil {
		return nil, err
	}

	return &DataBase{DB: pool}, nil
}

