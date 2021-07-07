package persistence

import (
	"context"
	"errors"

	"github.com/Timos-API/authenticator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserPersistor struct {
	c *mongo.Collection
}

var (
	ErrInsert = errors.New("error after inserting document")
)

func NewUserPersistor(c *mongo.Collection) *UserPersistor {
	return &UserPersistor{c}
}

func (p *UserPersistor) FindById(ctx context.Context, id string) (*authenticator.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	res := p.c.FindOne(ctx, bson.M{"_id": oid})

	if res.Err() != nil {
		return nil, res.Err()
	}

	var user authenticator.User
	err = res.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *UserPersistor) FindByProvider(ctx context.Context, provider string, providerId string) (*authenticator.User, error) {
	res := p.c.FindOne(ctx, bson.M{"providerId": providerId, "provider": provider})

	if res.Err() != nil {
		return nil, res.Err()
	}

	var user authenticator.User
	err := res.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *UserPersistor) CreateUser(ctx context.Context, user authenticator.User) (*authenticator.User, error) {

	res, err := p.c.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return p.FindById(ctx, oid.Hex())
	}

	return nil, ErrInsert
}

func (p *UserPersistor) UpdateUser(ctx context.Context, id string, update bson.M) (*authenticator.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	res := p.c.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if res.Err() != nil {
		return nil, res.Err()
	}

	var user authenticator.User
	err = res.Decode(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
