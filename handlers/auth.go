package handlers

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"thesis-management-app/pkgs/ldap"
	"thesis-management-app/pkgs/server"
	"thesis-management-app/views/auth"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, auth.Login())
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := ldap.UserCredentials{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}
	fmt.Printf("cred %v", credentials)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
