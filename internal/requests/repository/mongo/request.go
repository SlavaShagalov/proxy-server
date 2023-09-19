package mongo

import (
	"github.com/SlavaShagalov/proxy-server/internal/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Req  requests.Req       `bson:"request"`
	Resp requests.Resp      `bson:"response"`
}

func fromModel(r *requests.Request) *Request {
	return &Request{
		Req:  r.Req,
		Resp: r.Resp,
	}
}

func (r *Request) ToModel() *requests.Request {
	return &requests.Request{
		ID:   r.ID.Hex(),
		Req:  r.Req,
		Resp: r.Resp,
	}
}
