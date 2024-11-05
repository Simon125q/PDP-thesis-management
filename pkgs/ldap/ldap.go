package ldap

import (
	"fmt"
	"os"

	"github.com/go-ldap/ldap/v3"
)

type UserCredentials struct {
	Email    string
	Password string
}

type LDAPConfig struct {
	ServerURL    string
	BaseDN       string
	BindDN       string
	BindPassword string
}

func LoadLDAPConfig() LDAPConfig {
	return LDAPConfig{
		ServerURL:    os.Getenv("LDAP_SERVER_URL"),
		BaseDN:       os.Getenv("LDAP_BASE_DN"),
		BindDN:       os.Getenv("LDAP_BIND_DN"),
		BindPassword: os.Getenv("LDAP_BIND_PASSWORD"),
	}
}

func ConnectLDAP(cfg LDAPConfig) (*ldap.Conn, error) {
	l, err := ldap.DialURL(cfg.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP server: %v", err)
	}
	err = l.Bind(cfg.BindDN, cfg.BindPassword)
	if err != nil {
		l.Close()
		return nil, fmt.Errorf("LDAP bind failed: %v", err)
	}
	return l, nil
}
