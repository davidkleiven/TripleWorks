package api

import "context"

const defaultUser string = "Unknown"
const userEmail string = "userEmail"
const userRoleKey string = "userRole"

func UserFromCtx(ctx context.Context) string {
	user, _ := ctx.Value(userEmail).(string)
	if user == "" {
		return defaultUser
	}
	return user
}

func RoleFromCtx(ctx context.Context) string {
	role, _ := ctx.Value(userRoleKey).(string)
	return role
}
