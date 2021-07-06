package service

import (
	"Timos-API/Accounts/helper"
	"context"
	"os"
	"strings"
	"time"

	authenticator "github.com/Timos-API/Authenticator"
	"github.com/brianvoe/sjwt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitter"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthService struct {
	u *UserService
}

type JwtToken struct {
	Token string `json:"token"`
}

func NewAuthService(u *UserService) *AuthService {
	return &AuthService{u}
}

func (s *AuthService) RegisterOAuth() {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = strings.HasPrefix(os.Getenv("CALLBACK"), "https")

	gothic.Store = store

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.ExpandEnv("${CALLBACK}/auth/google/callback"), "profile"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.ExpandEnv("${CALLBACK}/auth/github/callback"), "user:name"),
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), os.ExpandEnv("${CALLBACK}/auth/twitter/callback")),
	)
}

func (s *AuthService) UserSignedIn(ctx context.Context, gothUser goth.User) (*JwtToken, error) {

	millis := helper.CurrentTimeMillis()
	var jwtUser *authenticator.User

	if exist, user := s.u.doesUserExistP(ctx, gothUser.Provider, gothUser.UserID); exist {

		user, err := s.u.updateUser(ctx, user.UserID.Hex(), bson.M{
			"avatar":     gothUser.AvatarURL,
			"last_login": millis,
			"name":       gothUser.Name,
		})

		if err != nil {
			return nil, err
		}
		jwtUser = user

	} else {

		user, err := s.u.createUser(ctx, authenticator.User{
			ProviderID:  gothUser.UserID,
			Provider:    gothUser.Provider,
			Name:        gothUser.Name,
			Avatar:      gothUser.AvatarURL,
			Group:       "user",
			MemberSince: millis,
			LastLogin:   millis,
		})

		if err != nil {
			return nil, err
		}
		jwtUser = user

	}

	claims, err := sjwt.ToClaims(jwtUser)
	if err != nil {
		return nil, err
	}

	claims.SetExpiresAt(time.Now().Add(time.Hour * 24))
	jwt := claims.Generate([]byte(os.Getenv("JWT_SECRET")))

	return &JwtToken{jwt}, nil

}
