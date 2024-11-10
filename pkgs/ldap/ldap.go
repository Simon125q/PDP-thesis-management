package ldap

import (
	"context"
	"os"
	"thesis-management-app/pkgs/server"
	"time"

	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/ldap"
	"github.com/shaj13/go-guardian/store"
)

type UserCredentials struct {
	Login    string
	Password string
}

func LoadLDAPConfig() *ldap.Config {
	return &ldap.Config{
		BaseDN:       os.Getenv("LDAP_BASE_DN"),
		BindDN:       os.Getenv("LDAP_BIND_DN"),
		Port:         os.Getenv("LDAP_PORT"),
		Host:         os.Getenv("LDAP_HOST"),
		BindPassword: os.Getenv("LDAP_BIND_PASSWORD"),
		Filter:       os.Getenv("LDAP_FILTER"),
	}
}

func SetupLDAP() {
	cfg := LoadLDAPConfig()
	server.MyS.Authenticator = auth.New()
	server.MyS.Cache = store.NewFIFO(context.Background(), time.Minute*10)
	strategy := ldap.NewCached(cfg, server.MyS.Cache)
	server.MyS.Authenticator.EnableStrategy(ldap.StrategyKey, strategy)
}

// func MockLDAPAuthenticate(cred UserCredentials) (LDAPResponse, error) {
// 	user1 := UserCredentials{
// 		Email:    "admin@gmail.com",
// 		Password: "admin123",
// 	}
// 	user2 := UserCredentials{
// 		Email:    "user@gmail.com",
// 		Password: "user123",
// 	}
// 	if cred.Email == user1.Email && cred.Password == user1.Password {
// 		return LDAPResponse{
// 			Email:       user1.Email,
// 			Name:        "admin1",
// 			AccessToken: "ldap1",
// 		}, nil
// 	}
// 	if cred.Email == user2.Email && cred.Password == user2.Password {
// 		return LDAPResponse{
// 			Email:       user2.Email,
// 			Name:        "user1",
// 			AccessToken: "ldap2",
// 		}, nil
// 	}
// 	return LDAPResponse{}, fmt.Errorf("Wrong credentials")
// }
