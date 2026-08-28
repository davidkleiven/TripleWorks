package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const jwtCookieName = "tripleworksJwt"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUnauthorized = errors.New("unauthorized")
)

type Auth struct {
	ClientId      string
	ClientSecret  string
	SessionSecret string
	Callback      string
	JwtSecret     []byte
	JwtTtl        time.Duration
}

type UserStore interface {
	GetByEmail(ctx context.Context, email string) (models.User, error)
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

func MakeHandleAuthCallback(
	userAuth func(w http.ResponseWriter, r *http.Request) (goth.User, error),
	store UserStore,
	jwtSecret []byte,
	jwtTtl time.Duration,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userAuth(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := store.GetByEmail(r.Context(), user.Email)
		if err != nil {
			slog.Info("Login rejected: user not allowed", "email", user.Email)
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		token, err := NewToken(jwtSecret, dbUser, jwtTtl)
		if err != nil {
			http.Error(w, "Could not create token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     jwtCookieName,
			Value:    token,
			Path:     "/",
			MaxAge:   int(jwtTtl.Seconds()),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		slog.Info("Successful login", "provider", user.Provider, "email", dbUser.Email, "role", dbUser.Role)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func isPublicPath(p string) bool {
	return strings.HasPrefix(p, "/auth/") || strings.HasPrefix(p, "/js/")
}

func RequireAuth(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if isPublicPath(r.URL.Path) {
					next.ServeHTTP(w, r)
					return
				}

				token, err := r.Cookie(jwtCookieName)
				if err != nil {
					http.Redirect(w, r, "/auth/google", http.StatusSeeOther)
					return
				}

				claims, err := ParseToken(jwtSecret, token.Value)
				if err != nil {
					http.Redirect(w, r, "/auth/google", http.StatusSeeOther)
					return
				}

				ctx := context.WithValue(r.Context(), userEmail, claims.Email)
				ctx = context.WithValue(ctx, userRoleKey, claims.Role)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
	}
}
