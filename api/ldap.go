package api

import (
	"sort"

	"github.com/korylprince/easyldap"
	"github.com/korylprince/go-ad-auth"
)

//LDAPDB represents an LDAP database
type LDAPDB interface {
	List() ([]*LDAPUser, error)
}

//ADDB represents an Active Directory database
type ADDB struct {
	config *easyldap.Config
}

//NewADDB creates a new ADDB
func NewADDB(bindDN, bindPass string, config *auth.Config) *ADDB {
	c := &easyldap.Config{
		Server:     config.Server,
		Port:       config.Port,
		BaseDN:     config.BaseDN,
		Security:   easyldap.SecurityType(config.Security),
		TLSConfig:  config.TLSConfig,
		PagingSize: 1000,
		Filter:     "(&(objectCategory=Person)(memberOf=CN=Staff,OU=User Groups,OU=Accounts,OU=bisd,DC=bullardisd,DC=net))",
		Attributes: []string{"sn", "givenname", "sAMAccountName", "mail"},
		Username:   bindDN,
		Password:   bindPass,
	}
	return &ADDB{config: c}
}

//List gets and returns a new list of Users or an error if one occurred
func (db *ADDB) List() ([]*LDAPUser, error) {
	entries, err := easyldap.Query(db.config)
	if err != nil {
		return nil, err
	}
	var s []*LDAPUser
	for _, e := range entries {
		user := &LDAPUser{
			LastName:  e.GetAttributeValue("sn"),
			FirstName: e.GetAttributeValue("givenName"),
			Username:  e.GetAttributeValue("sAMAccountName"),
			Email:     e.GetAttributeValue("mail"),
		}
		s = append(s, user)
	}
	sort.Sort(UserSorter(s))
	return s, nil
}

//LDAPUser represents a User in an LDAP Database
type LDAPUser struct {
	LastName  string
	FirstName string
	Username  string
	Email     string
}

func (u LDAPUser) sortKey() string {
	return u.LastName + u.FirstName
}

//UserSorter is a helper type to sort Users
type UserSorter []*LDAPUser

//Len is a helper function to sort Users
func (u UserSorter) Len() int {
	return len(u)
}

//Swap is a helper function to sort Users
func (u UserSorter) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

//Less is a helper function to sort Users
func (u UserSorter) Less(i, j int) bool {
	return u[i].sortKey() < u[j].sortKey()
}
