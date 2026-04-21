package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const sessionKey = "tripleworksSession"

type Auth struct {
	ClientId      string
	ClientSecret  string
	SessionSecret string
	Callback      string
}

func (a *Auth) EnsureInitialized() {
	data := map[string]string{
		"clientId":      a.ClientId,
		"clientSecret":  a.ClientSecret,
		"sessionSecret": a.SessionSecret,
		"callback":      a.Callback,
	}

	for k, v := range data {
		if v == "" {
			panic(fmt.Sprintf("%s is empty", k))
		}
	}
}

func (a *Auth) Setup() {
	goth.UseProviders(
		google.New(a.ClientId, a.ClientSecret, a.Callback, "email", "profile"),
	)

	key := []byte(a.SessionSecret)
	store := sessions.NewCookieStore(key)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
	}
	gothic.Store = store
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func MakeHandleAuthCallback(userAuth func(w http.ResponseWriter, r *http.Request) (goth.User, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userAuth(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		session, err := gothic.Store.Get(r, sessionKey)
		if err != nil {
			http.Error(w, "Could not fetch session: "+err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["userEmail"] = user.Email
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Could not store session: "+err.Error(), http.StatusInternalServerError)
			return
		}
		slog.Info("Successful login", "provider", user.Provider)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func GetUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			session, err := gothic.Store.Get(r, sessionKey)
			if err != nil {
				http.Error(w, "Failed to get session", http.StatusInternalServerError)
				return
			}

			userValue := session.Values["userEmail"]
			if userValue == nil {
				http.Redirect(w, r, "/auth/google", http.StatusSeeOther)
				return
			}

			user := userValue.(string)
			ctx := context.WithValue(r.Context(), userEmail, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}
