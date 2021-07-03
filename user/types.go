package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserInfo struct {
	UserID primitive.ObjectID `json:"id"`
	Name   string             `json:"name"`
	Avatar string             `json:"avatar"`
}
