package db

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/SlavaShagalov/proxy-server/internal/pkg/config"
)

func NewMongoDB(log *zap.Logger) (*mongo.Database, error) {
	log.Info("Connecting to MongoDB...",
		zap.String("host", viper.GetString(config.MongoHost)),
		zap.Int("port", viper.GetInt(config.MongoPort)),
		zap.String("database", viper.GetString(config.MongoDB)),
	)

	//uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
	//	viper.GetString(config.MongoUser),
	//	viper.GetString(config.MongoPassword),
	//	viper.GetString(config.MongoHost),
	//	viper.GetInt(config.MongoPort),
	//	viper.GetString(config.MongoDB),
	//)

	uri := fmt.Sprintf("mongodb://%s:%d/%s",
		viper.GetString(config.MongoHost),
		viper.GetInt(config.MongoPort),
		viper.GetString(config.MongoDB),
	)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Error("Failed to create MongoDB connection", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create MongoDB connection")
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Error("Failed to connect to MongoDB", zap.Error(err))
		return nil, errors.Wrap(err, "failed to connect to MongoDB")
	}

	db := client.Database(viper.GetString(config.MongoDB))

	log.Info("MongoDB connection created successfully")
	return db, nil
}
