package datasource

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	*pgx.Conn
}

func NewDB(dsn string) *DB {
	return &DB{
		newPgConn(dsn),
	}
}

// Create new postgres connection, don't forget to close connection
// after using it
func newPgConn(dsn string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Errorf("failed to create database connection\n%v", err))
	}

	return conn
}
