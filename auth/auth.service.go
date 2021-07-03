package auth

import (
	"Timos-API/Accounts/user"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func RegisterOAuth() {

	fmt.Println("Register OAuth")

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = strings.HasPrefix(os.Getenv("CALLBACK"), "https")

	gothic.Store = store

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.ExpandEnv("${CALLBACK}/auth/google/callback"), "profile"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.ExpandEnv("${CALLBACK}/auth/github/callback"), "user:name"),
	)
}

func handleOAuthCallback(w http.ResponseWriter, req *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, req)

	fmt.Printf("Goth User %v\n", gothUser)

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	user := user.UserSignedIn(gothUser)

	if user == nil {
		fmt.Fprintln(w, errors.New("USER NOT FOUND"))
		return
	}

	claims, _ := sjwt.ToClaims(user)
	claims.SetExpiresAt(time.Now().Add(time.Hour * 24))
	jwt := claims.Generate([]byte(os.Getenv("JWT_SECRET")))

	t, err := template.ParseFiles("auth.html")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	t.Execute(w, JwtToken{jwt})

}
