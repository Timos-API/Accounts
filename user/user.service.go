package user

import (
	"Timos-API/Accounts/database"
	ctx "context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	authenticator "github.com/Timos-API/Authenticator"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func collection() *mongo.Collection {
	return database.Database.Collection("user")
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := GetUserById(params["id"])

	if user == nil {
		json.NewEncoder(w).Encode(authenticator.Exception{Message: "User not found"})
		return
	}

	json.NewEncoder(w).Encode(UserInfo{user.UserID, user.Name, user.Avatar})
}

func GetUserById(userId string) *authenticator.User {
	var user authenticator.User
	oid, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		fmt.Printf("Invalid ObjectID %v\n", userId)
		return nil
	}

	err = collection().FindOne(ctx.Background(), bson.M{"_id": oid}).Decode(&user)

	if err != nil {
		fmt.Printf("User not found... (%v) %v \n", userId, err)
		return nil
	}

	return &user
}

func UserSignedIn(gothUser goth.User) *authenticator.User {
	user := getUserByProvider(gothUser.UserID, gothUser.Provider)

	if user == nil {
		return registerUser(gothUser)
	} else {
		return updateUser(user, gothUser)
	}
}

func getUserByProvider(providerId string, provider string) *authenticator.User {
	var user authenticator.User
	err := collection().FindOne(ctx.Background(), bson.M{"providerId": providerId, "provider": provider}).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func registerUser(gothUser goth.User) *authenticator.User {
	t := time.Now()
	tUnixMilli := int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)

	var user = authenticator.User{
		ProviderID:  gothUser.UserID,
		Provider:    gothUser.Provider,
		Name:        gothUser.Name,
		Avatar:      gothUser.AvatarURL,
		Group:       "user",
		MemberSince: tUnixMilli,
		LastLogin:   tUnixMilli,
	}

	fmt.Println("Register User")

	result, _ := collection().InsertOne(ctx.Background(), user)

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return GetUserById(oid.Hex())
	}

	fmt.Println("Returning nil...")
	return nil
}

func updateUser(user *authenticator.User, gothUser goth.User) *authenticator.User {
	t := time.Now()
	tUnixMilli := int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)

	fmt.Println("Update User")

	collection().UpdateByID(ctx.Background(), user.UserID, bson.M{"$set": bson.M{
		"avatar":     gothUser.AvatarURL,
		"last_login": tUnixMilli,
		"name":       gothUser.Name,
	}})

	return GetUserById(user.UserID.Hex())
}
