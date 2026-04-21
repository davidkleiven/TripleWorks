package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/testify/require"
)

func TestSetupWithEmptyOk(t *testing.T) {
	auth := Auth{}
	require.NotPanics(t, func() { auth.Setup() })
}

type FailingSessionStore struct {
	GetErr  error
	SaveErr error
	NewErr  error
	Session *sessions.Session
}

func (f *FailingSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return f.Session, f.GetErr
}

func (f *FailingSessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return f.SaveErr
}

func (f *FailingSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return nil, f.NewErr
}

func createAuthenticatedRequest() *http.Request {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	session, _ := gothic.Store.Get(req, sessionKey)
	if session == nil {
		return req
	}
	session.Values["userEmail"] = "admin@example.com"
	session.Save(req, rec)

	// Get the cookie from the recorder
	res := rec.Result()
	cookie := res.Cookies()[0]
	req.AddCookie(cookie)
	return req
}

func TestGetUserMiddleware(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	wrapped := GetUserMiddleware(http.HandlerFunc(handler))

	gothic.Store = &FailingSessionStore{GetErr: errors.New("something went wrong")}

	t.Run("internal server error on no session cookie", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		wrapped.ServeHTTP(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	auth := Auth{SessionSecret: "top-secret"}
	auth.Setup()

	t.Run("redirect on missing user", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		wrapped.ServeHTTP(rec, req)
		require.Equal(t, http.StatusSeeOther, rec.Code)
	})

	t.Run("ok on valid session", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := createAuthenticatedRequest()
		wrapped.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestAuthPanicsOnEmpty(t *testing.T) {
	auth := Auth{}
	require.Panics(t, func() { auth.EnsureInitialized() })
}

func TestMakeHandleAuthCallback(t *testing.T) {
	auth := Auth{SessionSecret: "top-secret"}
	auth.Setup()

	validUser := func(w http.ResponseWriter, r *http.Request) (goth.User, error) {
		return goth.User{
			Email:    "test@example.com",
			Provider: "google",
		}, nil
	}

	t.Run("successful login redirects to home", func(t *testing.T) {
		handler := MakeHandleAuthCallback(validUser)

		mux := http.NewServeMux()
		mux.HandleFunc("/auth/callback", handler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/callback?code=test&state=state", nil)
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code)
		require.Equal(t, "/", rec.Header().Get("Location"))
	})

	t.Run("bad request on user auth error", func(t *testing.T) {
		mockUserAuth := func(w http.ResponseWriter, r *http.Request) (goth.User, error) {
			return goth.User{}, errors.New("auth failed")
		}
		handler := MakeHandleAuthCallback(mockUserAuth)

		mux := http.NewServeMux()
		mux.HandleFunc("/auth/callback", handler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/callback", nil)
		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "auth failed")
	})

	// Run test associated with session failure
	failingStore := FailingSessionStore{}
	gothic.Store = &failingStore

	handler := MakeHandleAuthCallback(validUser)
	t.Run("internal server error on store get error", func(t *testing.T) {
		defer func() {
			failingStore.GetErr = nil
		}()
		failingStore.GetErr = errors.New("something went wrong")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/callback", nil)
		handler(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "Could not fetch")
	})

	t.Run("internal server error on save failure", func(t *testing.T) {
		defer func() {
			failingStore.SaveErr = nil
			failingStore.Session = nil
		}()

		req := httptest.NewRequest("GET", "/auth", nil)

		failingStore.Session = sessions.NewSession(&failingStore, sessionKey)
		failingStore.SaveErr = errors.New("something went wrong")

		rec := httptest.NewRecorder()
		handler(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "Could not store")
	})
}
