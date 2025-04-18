package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/pkgs/sessions"
	"thesis-management-app/types"
	"time"
)

func getUserIdFromUsername(name string) int {
	ids, err := server.MyS.DB.GetPersonIDByFullName(name)
	if err != nil {
		slog.Error("getUserIdFromUsername", "err", err)
		return 0
	}
	if len(ids) == 0 {
		slog.Error("getUserIdFromUsername - Coudnt find any employee with name: " + name)
		return 0
	}
	id, _ := strconv.Atoi(ids[0])
	slog.Info("getUserIdFromUsername", "Logged in as user with id", id)
	return id
}

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
		var uId int
		if currentSession.IsAdmin {
			uId = 1
		} else {
			uId = getUserIdFromUsername(currentSession.Username)
		}
		user := types.AuthenticatedUser{
			Id:       uId,
			Login:    currentSession.Username,
			IsAdmin:  currentSession.IsAdmin,
			LoggedIn: true,
		}
		slog.Info("WithUser", "User", user)
		slog.Info("WithUser", "URL", r.URL.Path)
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAutehenticatedUser(r)
		slog.Info("WithAuth", "url", r.URL.Path)
		if !user.LoggedIn {
			slog.Info("WithAuth", "login", false)
			hxRedirect(w, r, "/login")
			// http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if !user.IsAdmin {
			q := r.URL.Query()
			q.Add("user_id", strconv.Itoa(user.Id))
			r.URL.RawQuery = q.Encode()
			slog.Info("withAuth", "RawQuery", r.URL.RawQuery)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func WithAdminRights(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAutehenticatedUser(r)
		if !user.LoggedIn || !user.IsAdmin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func RefreshSession(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		currSession, ok := sessions.Sessions.Get(cookie.Value)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		if currSession.IsExpired() {
			newToken, _ := sessions.Sessions.Refresh(cookie.Value)
			http.SetCookie(w, &http.Cookie{
				Name:  "session_token",
				Value: newToken,
				Path:  "/",
				// Domain:  "10.1.1.180",
				Expires: time.Now().Add(120 * time.Minute),
				Secure:  false,
			})
		}
		next.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}
