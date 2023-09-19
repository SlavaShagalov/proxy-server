package mongo

import (
	"context"
	"errors"
	"github.com/SlavaShagalov/proxy-server/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/proxy-server/internal/pkg/errors"
	pRequests "github.com/SlavaShagalov/proxy-server/internal/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type repository struct {
	coll *mongo.Collection
	log  *zap.Logger
}

func New(coll *mongo.Collection, log *zap.Logger) pRequests.Repository {
	return &repository{
		coll: coll,
		log:  log,
	}
}

func (r *repository) Create(params *pRequests.Request) error {
	_, err := r.coll.InsertOne(context.TODO(), fromModel(params))
	return err
}

func (r *repository) Get(id string) (*pRequests.Request, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	mongoRequest := new(Request)
	err = r.coll.FindOne(context.TODO(), filter).Decode(mongoRequest)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, pErrors.ErrRequestNotFound
		}

		r.log.Error(constants.DBError, zap.Error(err))
		return nil, pErrors.ErrDb
	}

	return mongoRequest.ToModel(), nil
}

func (r *repository) List() ([]pRequests.Request, error) {
	filter := bson.M{}
	cur, err := r.coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	requests := []pRequests.Request{}

	for cur.Next(context.Background()) {
		mongoRequest := new(Request)
		err = cur.Decode(mongoRequest)
		if err != nil {
			r.log.Error(constants.DBError, zap.Error(err))
			return nil, pErrors.ErrDb
		}
		requests = append(requests, *mongoRequest.ToModel())
	}

	return requests, nil
}
