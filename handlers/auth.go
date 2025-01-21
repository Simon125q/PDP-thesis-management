package handlers

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/ldap"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/pkgs/sessions"
	"thesis-management-app/views/auth"
	"time"

	"github.com/google/uuid"
	goAuth "github.com/shaj13/go-guardian/auth"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, auth.Login())
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := ldap.UserCredentials{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials.Login+":"+credentials.Password))
	r.Header.Set("Authorization", authHeader)
	user, err := server.MyS.Authenticator.Authenticate(r)
	if err != nil {
		slog.Error("coudnt authenticate", "err", err)
		loginErrs := auth.LoginErrors{
			InvalidCredentials: fmt.Sprintf("coudnt authenticate user, error occurred: %v", err),
		}
		return Render(w, r, auth.LoginForm(credentials, loginErrs))
	}
	//slog.Info("Authenticated", "user", user)

	sessionToken := uuid.NewString()
	newSession := sessions.Session{
		Username: credentials.Login,
		IsAdmin:  isAdmin(user),
		Expiry:   time.Now().Add(480 * time.Second),
	}
	sessions.Sessions.Add(sessionToken, newSession)
	setAuthCookie(w, sessionToken)

	return hxRedirect(w, r, "/")
}

func HandleLogoutPost(w http.ResponseWriter, r *http.Request) error {
	currCookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	sessions.Sessions.Remove(currCookie.Value)
	setAuthCookie(w, "")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func setAuthCookie(w http.ResponseWriter, sessionToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: time.Now().Add(120 * time.Minute),
		Secure:  false,
	})
}

func isAdmin(user goAuth.Info) bool {
	// groups := user.Groups()
	// for _, group := range groups {
	//    if group == "admin"{
	//         return true
	//    }
	// }
	if user.UserName() == "tesla" {
		return true
	}
	return false
}
