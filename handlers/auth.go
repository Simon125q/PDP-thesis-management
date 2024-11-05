package handlers

import (
	"fmt"
	"net/http"
	"thesis-management-app/pkgs/ldap"
	"thesis-management-app/views/auth"
	"thesis-management-app/views/home"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, auth.Login())
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := ldap.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	resp, err := ldap.MockLDAPAuthenticate(credentials)
	if err != nil {
		loginErrs := auth.LoginErrors{
			InvalidCredentials: "some error occured",
		}
		return Render(w, r, auth.LoginForm(credentials, loginErrs))
	}

	fmt.Printf("%v\n", resp)

	return Render(w, r, home.Index())
}
