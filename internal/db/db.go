package db

import (
	"context"
	"github.com/defer-panic/url-shortener-api/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func Connect(ctx context.Context, config config.DBConfig) (*DB, error) {
	uri := options.Client().ApplyURI(config.DSN)
	if config.User != "" && config.Password != "" {
		uri.SetAuth(options.Credential{Username: config.User, Password: config.Password, AuthSource: config.Database})
	}
	client, err := mongo.Connect(ctx, uri)
	if err != nil {
		return nil, err
	}

	return &DB{client: client}, nil
}

func (d *DB) Client() *mongo.Client {
	return d.client
}

func (d *DB) Close(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}
