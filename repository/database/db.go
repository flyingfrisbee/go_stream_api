package database

import (
	"context"
	"go_stream_api/repository/database/relational"
)

type persistentStorageType int

const (
	Postgres persistentStorageType = iota
	MongoDB
)

type dbConn struct {
	cancel context.CancelFunc
	*relational.PostgresInstance
}

var (
	Conn dbConn
)

func StartConnectionToDB() {
	ctx, cancel := context.WithCancel(context.Background())
	Conn = dbConn{
		cancel:           cancel,
		PostgresInstance: relational.StartConnection(ctx),
	}
}

func TerminateConnectionToDB() {
	Conn.ClosePostgresConnection()
	Conn.cancel()
}
