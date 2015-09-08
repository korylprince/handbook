package main

import (
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/korylprince/go-ad-auth"
)

//Config represents options given in the environment
type Config struct {
	LDAPServer     string //required
	LDAPPort       int    //default: 389
	LDAPBaseDN     string //required
	LDAPGroup      string //optional
	LDAPAdminGroup string //optional
	LDAPSecurity   string //default: none
	ldapSecurity   auth.SecurityType

	LDAPBindDN       string //required, can be in user@dmain.tld format
	LDAPBindPass     string //required
	LDAPSearchBaseDN string //required

	SessionDuration      int //in minutes; default: 5
	AdminSessionDuration int //in minutes; default: 60

	SQLDriver string //required
	SQLDSN    string //required

	ListenAddr string //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash

	Debug bool //default: false
}

var config = &Config{}

func checkEmpty(val, name string) {
	if val == "" {
		log.Fatalf("HANDBOOK_%s must be configured\n", name)
	}
}

func init() {
	err := envconfig.Process("HANDBOOK", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}
	checkEmpty(config.LDAPServer, "LDAPSERVER")

	if config.LDAPPort == 0 {
		config.LDAPPort = 389
	}

	checkEmpty(config.LDAPBaseDN, "LDAPBASEDN")

	switch strings.ToLower(config.LDAPSecurity) {
	case "", "none":
		config.ldapSecurity = auth.SecurityNone
	case "tls":
		config.ldapSecurity = auth.SecurityTLS
	case "starttls":
		config.ldapSecurity = auth.SecurityStartTLS
	default:
		log.Fatalln("Invalid HANDBOOK_LDAPSECURITY:", config.LDAPSecurity)
	}

	checkEmpty(config.LDAPBindDN, "LDAPBINDDN")
	checkEmpty(config.LDAPBindPass, "LDAPBINDPASS")
	checkEmpty(config.LDAPSearchBaseDN, "LDAPSEARCHBASEDN")

	if config.SessionDuration == 0 {
		config.SessionDuration = 5
	}

	if config.AdminSessionDuration == 0 {
		config.AdminSessionDuration = 60
	}

	checkEmpty(config.SQLDriver, "SQLDriver")
	checkEmpty(config.SQLDSN, "SQLDSN")

	if config.SQLDriver == "mysql" && !strings.Contains(config.SQLDSN, "?parseTime=true") {
		log.Fatalln("mysql DSN must contain \"?parseTime=true\"")
	}

	checkEmpty(config.ListenAddr, "LISTENADDR")
}
