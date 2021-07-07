package service

import (
	"Timos-API/Accounts/persistence"
	"context"
	"errors"

	"github.com/Timos-API/authenticator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	p *persistence.UserPersistor
}

type UserInfo struct {
	UserID      primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Avatar      string             `json:"avatar"`
	MemberSince int64              `json:"member_since" bson:"member_since"`
}

var (
	ErrInvalidObjectID = errors.New("invalid ObjectID")
)

func NewUserService(p *persistence.UserPersistor) *UserService {
	return &UserService{p}
}

func (s *UserService) GetUserInfo(ctx context.Context, id string) (*UserInfo, error) {
	user, err := s.p.FindById(ctx, id)

	if err != nil {
		return nil, err
	}
	return &UserInfo{UserID: user.UserID, Name: user.Name, Avatar: user.Avatar, MemberSince: user.MemberSince}, nil
}

func (s *UserService) createUser(ctx context.Context, user authenticator.User) (*authenticator.User, error) {
	return s.p.CreateUser(ctx, user)
}

func (s *UserService) updateUser(ctx context.Context, id string, update bson.M) (*authenticator.User, error) {
	return s.p.UpdateUser(ctx, id, update)
}

func (s *UserService) doesUserExist(ctx context.Context, provider string, providerId string) (bool, *authenticator.User) {
	user, err := s.p.FindByProvider(ctx, provider, providerId)
	return err == nil && user != nil, user
}
