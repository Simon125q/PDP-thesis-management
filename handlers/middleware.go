package handlers

import (
	"context"
	"net/http"
	"strings"
	"thesis-management-app/pkgs/sessions"
	"thesis-management-app/types"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		currentSession, ok := sessions.Sessions.Get(cookie.Value)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		user := types.AuthenticatedUser{
			Login:    currentSession.Username,
			IsAdmin:  currentSession.IsAdmin,
			LoggedIn: true,
		}
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
