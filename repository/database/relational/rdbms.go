package relational

import (
	"context"
	env "go_stream_api/environment"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresInstance struct {
	postgresConn
	animeRelatedQuery
}

type postgresConn struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func (conn *postgresConn) initializeTables() {
	sqlBytes, err := os.ReadFile("./repository/database/relational/init.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.pool.Exec(conn.ctx, string(sqlBytes))
	if err != nil {
		log.Fatal(err)
	}
}

func (conn *postgresConn) CloseConnection() {
	conn.pool.Close()
	conn.pool = nil
}

func StartConnection(ctx context.Context) *PostgresInstance {
	pool, err := pgxpool.New(ctx, env.RDBMSConnString)
	if err != nil {
		log.Fatal(err)
	}

	pgConn := postgresConn{
		ctx:  ctx,
		pool: pool,
	}

	instance := &PostgresInstance{
		postgresConn: pgConn,
		animeRelatedQuery: &animeTable{
			conn: pgConn,
		},
	}

	instance.initializeTables()

	return instance
}
