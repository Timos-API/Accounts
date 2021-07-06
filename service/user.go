package service

import (
	"Timos-API/Accounts/persistence"
	"context"
	"errors"

	authenticator "github.com/Timos-API/Authenticator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	p *persistence.UserPersistor
}

type UserInfo struct {
	UserID primitive.ObjectID `json:"id"`
	Name   string             `json:"name"`
	Avatar string             `json:"avatar"`
}

var (
	ErrInvalidObjectID = errors.New("Invalid ObjectID")
)

func NewUserService(p *persistence.UserPersistor) *UserService {
	return &UserService{p}
}

func (s *UserService) getUserById(ctx context.Context, id string) (*authenticator.User, error) {
	if primitive.IsValidObjectID(id) {
		return nil, ErrInvalidObjectID
	}
	return s.p.FindById(ctx, id)
}

func (s *UserService) GetUserInfo(ctx context.Context, id string) (*UserInfo, error) {
	user, err := s.p.FindById(ctx, id)

	if err != nil {
		return nil, err
	}
	return &UserInfo{UserID: user.UserID, Name: user.Name, Avatar: user.Avatar}, nil
}

func (s *UserService) createUser(ctx context.Context, user authenticator.User) (*authenticator.User, error) {
	return s.p.CreateUser(ctx, user)
}

func (s *UserService) updateUser(ctx context.Context, id string, update authenticator.User) (*authenticator.User, error) {
	if primitive.IsValidObjectID(id) {
		return nil, ErrInvalidObjectID
	}
	return s.p.UpdateUser(ctx, id, update)
}

func (s *UserService) doesUserExist(ctx context.Context, id string) (bool, *authenticator.User) {
	if !primitive.IsValidObjectID(id) {
		return false, nil
	}
	user, err := s.getUserById(ctx, id)
	return err == nil && user != nil, user
}

func (s *UserService) doesUserExistP(ctx context.Context, provider string, providerId string) (bool, *authenticator.User) {
	user, err := s.p.FindByProvider(ctx, provider, providerId)
	return err == nil && user != nil, user
}
