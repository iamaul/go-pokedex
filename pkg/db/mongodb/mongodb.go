package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/iamaul/go-pokedex/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 10 * time.Second

// NewClient established connection to a mongoDb instance using provided URI and auth credentials.
func NewClient(c *config.Config) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(c.MongoDB.MongoURI)
	if c.MongoDB.Username != "" && c.MongoDB.Password != "" {
		opts.SetAuth(options.Credential{
			Username: c.MongoDB.Username, Password: c.MongoDB.Password,
		})
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func IsDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}
