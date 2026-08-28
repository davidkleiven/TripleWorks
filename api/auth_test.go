package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
)

type fakeUserStore struct {
	users map[string]models.User
	err   error
}

func (f *fakeUserStore) GetByEmail(_ context.Context, email string) (models.User, error) {
	if f.err != nil {
		return models.User{}, f.err
	}
	u, ok := f.users[email]
	if !ok {
		return models.User{}, ErrUserNotFound
	}
	return u, nil
}

func TestSetupWithEmptyOk(t *testing.T) {
	auth := Auth{}
	require.NotPanics(t, func() { auth.Setup() })
}

func TestAuthPanicsOnEmpty(t *testing.T) {
	auth := Auth{}
	require.Panics(t, func() { auth.EnsureInitialized() })
}

func TestJWT(t *testing.T) {
	secret := []byte("top-secret")
	user := models.User{Email: "test@example.com", Role: models.RoleAdmin}

	token, err := NewToken(secret, user, time.Hour)
	require.NoError(t, err)

	t.Run("round trip", func(t *testing.T) {
		claims, err := ParseToken(secret, token)
		require.NoError(t, err)
		require.Equal(t, user.Email, claims.Email)
		require.Equal(t, user.Role, claims.Role)
	})

	t.Run("expired token rejected", func(t *testing.T) {
		token, err := NewToken(secret, user, -time.Minute)
		require.NoError(t, err)
		_, err = ParseToken(secret, token)
		require.Error(t, err)
	})

	t.Run("tampered token rejected", func(t *testing.T) {
		_, err := ParseToken(secret, token+"tampered")
		require.Error(t, err)
	})

	t.Run("wrong secret rejected", func(t *testing.T) {
		_, err := ParseToken([]byte("other-secret"), token)
		require.Error(t, err)
	})
}

func jwtRequest(token string) *http.Request {
	req := httptest.NewRequest("GET", "/", nil)
	if token != "" {
		req.AddCookie(&http.Cookie{Name: jwtCookieName, Value: token})
	}
	return req
}

func TestRequireAuth(t *testing.T) {
	secret := []byte("top-secret")
	user := models.User{Email: "test@example.com", Role: models.RoleAdmin}
	token, err := NewToken(secret, user, time.Hour)
	require.NoError(t, err)

	var gotEmail, gotRole string
	handler := func(w http.ResponseWriter, r *http.Request) {
		gotEmail = UserFromCtx(r.Context())
		gotRole = RoleFromCtx(r.Context())
	}
	authMux := RequireAuth(secret)(http.HandlerFunc(handler))

	t.Run("redirect on missing cookie", func(t *testing.T) {
		rec := httptest.NewRecorder()
		authMux.ServeHTTP(rec, jwtRequest(""))
		require.Equal(t, http.StatusSeeOther, rec.Code)
		require.Equal(t, "/auth/google", rec.Header().Get("Location"))
	})

	t.Run("redirect on invalid token", func(t *testing.T) {
		rec := httptest.NewRecorder()
		authMux.ServeHTTP(rec, jwtRequest("not-a-token"))
		require.Equal(t, http.StatusSeeOther, rec.Code)
	})

	t.Run("handler runs with claims on valid token", func(t *testing.T) {
		rec := httptest.NewRecorder()
		authMux.ServeHTTP(rec, jwtRequest(token))
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, user.Email, gotEmail)
		require.Equal(t, user.Role, gotRole)
	})

	t.Run("public auth paths bypass", func(t *testing.T) {
		for _, p := range []string{"/auth/google", "/auth/google/callback", "/js/app.js"} {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", p, nil)
			authMux.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code, "path %s should be public", p)
		}
	})
}

func TestMakeHandleAuthCallback(t *testing.T) {
	secret := []byte("top-secret")
	ttl := time.Hour

	validGothUser := func(w http.ResponseWriter, r *http.Request) (goth.User, error) {
		return goth.User{Email: "test@example.com", Provider: "google"}, nil
	}

	store := &fakeUserStore{users: map[string]models.User{
		"test@example.com": {Email: "test@example.com", Role: models.RoleUser},
	}}

	t.Run("successful login sets jwt cookie", func(t *testing.T) {
		handler := MakeHandleAuthCallback(validGothUser, store, secret, ttl)

		mux := http.NewServeMux()
		mux.HandleFunc("/auth/callback", handler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/callback?code=test&state=state", nil)
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code)
		require.Equal(t, "/", rec.Header().Get("Location"))

		cookie := findCookie(rec.Result().Cookies(), jwtCookieName)
		require.NotNil(t, cookie)
		require.True(t, cookie.HttpOnly)
		require.Equal(t, "/", cookie.Path)

		claims, err := ParseToken(secret, cookie.Value)
		require.NoError(t, err)
		require.Equal(t, "test@example.com", claims.Email)
		require.Equal(t, models.RoleUser, claims.Role)
	})

	t.Run("bad request on user auth error", func(t *testing.T) {
		mockUserAuth := func(w http.ResponseWriter, r *http.Request) (goth.User, error) {
			return goth.User{}, errors.New("auth failed")
		}
		handler := MakeHandleAuthCallback(mockUserAuth, store, secret, ttl)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/callback", nil)
		handler(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "auth failed")
	})

	t.Run("unauthorized when user not allowed", func(t *testing.T) {
		notAllowed := func(w http.ResponseWriter, r *http.Request) (goth.User, error) {
			return goth.User{Email: "stranger@example.com", Provider: "google"}, nil
		}
		handler := MakeHandleAuthCallback(notAllowed, store, secret, ttl)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/callback", nil)
		handler(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
		require.Nil(t, findCookie(rec.Result().Cookies(), jwtCookieName))
	})
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}
