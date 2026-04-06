package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserFromContextDefaultOnMissing(t *testing.T) {
	user := UserFromCtx(context.Background())
	require.Equal(t, defaultUser, user)
}

type MockIdentifier struct {
	User string
	Err  error
}

func (m *MockIdentifier) WhoIs(ctx context.Context, addr string) (string, error) {
	return m.User, m.Err
}

func TestNoUserAddedOnError(t *testing.T) {
	identifier := MockIdentifier{Err: errors.New("something went wrong")}

	middleware := UserIdentificationMiddleware{
		Identifier: &identifier,
	}

	var user string
	handler := func(w http.ResponseWriter, r *http.Request) {
		user = UserFromCtx(r.Context())
	}
	wrappedHandler := middleware.Apply(http.HandlerFunc(handler))
	req := httptest.NewRequest("GET", "/endpoint", nil)
	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)
	require.Equal(t, defaultUser, user)

	identifier.Err = nil
	identifier.User = "username"
	wrappedHandler.ServeHTTP(rec, req)
	require.Equal(t, "username", user)
}

func TestTailscaleNotFound(t *testing.T) {
	tsIdentifier := TailscaleUserIdentifier{}
	user, err := tsIdentifier.WhoIs(context.Background(), "123.123")
	require.Error(t, err)
	require.Equal(t, "", user)
}
