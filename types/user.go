package types

const UserContextKey = "user"

type AuthenticatedUser struct {
	Login    string
	IsAdmin  bool
	LoggedIn bool
}
