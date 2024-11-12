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
