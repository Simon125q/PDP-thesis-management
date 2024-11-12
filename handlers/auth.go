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
	slog.Info("Authenticated", "user", user.UserName())

	sessionToken := uuid.NewString()
	newSession := sessions.Session{
		Username: credentials.Login,
		IsAdmin:  isAdmin(credentials.Login),
		Expiry:   time.Now().Add(480 * time.Second),
	}
	sessions.Sessions.Add(sessionToken, newSession)
	setAuthCookie(w, sessionToken)
	slog.Info("Cookie set", "sessionToken", sessionToken)
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	slog.Info("Cookie set", "cookie", cookie)
	slog.Info("Cookie set", "cookie val", cookie.Value)
	n, _ := sessions.Sessions.Get(cookie.Value)
	slog.Info("sessions", "store", n)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func setAuthCookie(w http.ResponseWriter, sessionToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(480 * time.Second),
	})
}

func isAdmin(username string) bool {
	if username == "tesla" {
		return true
	}
	return false
}

// if !validators.IsValidEmail(credentials.Login) {
// 	loginErrs := auth.LoginErrors{
// 		Email: "invalid email",
// 	}
// 	return Render(w, r, auth.LoginForm(credentials, loginErrs))
// }
// resp, err := ldap.MockLDAPAuthenticate(credentials)
// if err != nil {
// 	slog.Error("coudnt authenticate", "err", err)
// 	loginErrs := auth.LoginErrors{
// 		InvalidCredentials: fmt.Sprintf("coudnt authenticate user, error occurred: %v", err),
// 	}
// 	return Render(w, r, auth.LoginForm(credentials, loginErrs))
// }

// fmt.Printf("%v\n", resp)

// cookie := &http.Cookie{
// 	Value:    resp.AccessToken,
// 	Name:     "at",
// 	Path:     "/",
// 	HttpOnly: true,
// 	Secure:   true,
// }
// http.SetCookie(w, cookie)
