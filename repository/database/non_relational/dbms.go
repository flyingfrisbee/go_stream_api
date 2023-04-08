package nonrelational

import (
	"context"
	env "go_stream_api/environment"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBInstance struct {
	mongoDBConn
	episodeRelatedQuery
}

type mongoDBConn struct {
	ctx    context.Context
	client *mongo.Client
}

func (conn *mongoDBConn) CloseConnection() error {
	err := conn.client.Disconnect(conn.ctx)
	if err != nil {
		return err
	}

	conn.client = nil
	return nil
}

func StartConnection(ctx context.Context) *MongoDBInstance {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.DBMSConnString))
	if err != nil {
		log.Fatal(err)
	}

	mongoConn := mongoDBConn{
		ctx:    ctx,
		client: client,
	}

	return &MongoDBInstance{
		mongoDBConn: mongoConn,
		episodeRelatedQuery: &episodeCollections{
			conn: mongoConn,
		},
	}
}
