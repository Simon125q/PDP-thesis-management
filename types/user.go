package types

const UserContextKey = "user"

type AuthenticatedUser struct {
	Id       int
	Login    string
	IsAdmin  bool
	LoggedIn bool
}
