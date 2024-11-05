package handlers

import (
	"fmt"
	"net/http"
	"thesis-management-app/pkgs/ldap"
	"thesis-management-app/views/auth"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, auth.Login())
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := ldap.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	//check credentials in ldap

	loginErrs := auth.LoginErrors{
		InvalidCredentials: "some error occured",
	}
	fmt.Println(credentials)
	return Render(w, r, auth.LoginForm(credentials, loginErrs))
}
