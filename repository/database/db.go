package database

import (
	"context"
	nonrelational "go_stream_api/repository/database/non_relational"
	"go_stream_api/repository/database/relational"
	"log"
)

type PersistentStorageType int

const (
	Postgres PersistentStorageType = iota
	MongoDB
)

type dbConn struct {
	cancel context.CancelFunc
	Pg     *relational.PostgresInstance
	Mongo  *nonrelational.MongoDBInstance
}

var (
	Conn dbConn
)

func StartConnectionToDB() {
	ctx, cancel := context.WithCancel(context.Background())
	Conn = dbConn{
		cancel: cancel,
		Pg:     relational.StartConnection(ctx),
		Mongo:  nonrelational.StartConnection(ctx),
	}
}

func TerminateConnectionToDB() {
	Conn.Pg.CloseConnection()
	err := Conn.Mongo.CloseConnection()
	if err != nil {
		log.Fatal(err)
	}
	Conn.cancel()
}
