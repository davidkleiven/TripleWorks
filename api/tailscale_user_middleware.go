package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"tailscale.com/client/local"
)

type tailscaleUser string

const tailscaleUserKey tailscaleUser = "tailscaleUser"
const defaultUser string = "Unknown"

type UserIdentifier interface {
	WhoIs(ctx context.Context, addr string) (string, error)
}

type TailscaleUserIdentifier struct {
	Client local.Client
}

func (t *TailscaleUserIdentifier) WhoIs(ctx context.Context, addr string) (string, error) {
	who, err := t.Client.WhoIs(ctx, addr)
	if err != nil {
		return "", fmt.Errorf("Could not get user from tailscale: %w", err)
	}
	return who.UserProfile.LoginName, nil
}

type UserIdentificationMiddleware struct {
	Identifier UserIdentifier
}

func (t *UserIdentificationMiddleware) Apply(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			user, err := t.Identifier.WhoIs(r.Context(), r.RemoteAddr)
			if err != nil {
				slog.Warn("Could not derive user from tailscale", "error", err)
				h.ServeHTTP(w, r)
			} else {
				ctx := context.WithValue(r.Context(), tailscaleUserKey, user)
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		})
}

func UserFromCtx(ctx context.Context) string {
	value, ok := ctx.Value(tailscaleUserKey).(string)
	if !ok {
		return defaultUser
	}
	return value
}
