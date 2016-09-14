package api

import "github.com/korylprince/go-ad-auth"

//User represents a user's name
type User struct {
	EmployeeID string
	Username   string
	FirstName  string
	LastName   string
	Admin      bool
}

//Auth is an interface for an arbitrary authentication backend.
type Auth interface {
	//Login returns whether or not the given username or password is valid.
	//If valid, user will be non-nil
	//If the backend malfunctions, user will be nil and error will be non-nil.
	Login(username, password string) (user *User, err error)

	//AdminLogin returns whether or not the given username or password is valid admin login.
	//If valid, user will be non-nil
	//If the backend malfunctions, user will be nil and error will be non-nil.
	AdminLogin(username, password string) (user *User, err error)
}

//LDAPAuth represents an Auth that uses an Active Directory backend
type LDAPAuth struct {
	group      string
	adminGroup string
	config     *auth.Config
}

//NewLDAPAuth returns a new LDAPAuth with the given config, restricting logins to group and admins to to adminGroup if non-empty.
func NewLDAPAuth(group, adminGroup string, config *auth.Config) *LDAPAuth {
	return &LDAPAuth{
		group:      group,
		adminGroup: adminGroup,
		config:     config,
	}
}

//Login returns whether or not the given username or password is valid.
//If valid, user will be non-nil
//If the backend malfunctions, user will be nil and error will be non-nil.
func (a *LDAPAuth) Login(username, password string) (user *User, err error) {
	ok, attrs, err := auth.LoginWithAttrs(username, password, a.group, a.config, []string{"employeeID", "givenName", "sn"})
	if !ok {
		return nil, err
	}
	u := &User{Username: username}

	if ids, ok := attrs["employeeID"]; ok {
		if len(ids) > 0 {
			u.EmployeeID = ids[0]
		}
	}

	if gns, ok := attrs["givenName"]; ok {
		if len(gns) > 0 {
			u.FirstName = gns[0]
		}
	}

	if sns, ok := attrs["sn"]; ok {
		if len(sns) > 0 {
			u.LastName = sns[0]
		}
	}

	return u, err
}

//AdminLogin returns whether or not the given username or password is valid admin login.
//If valid, user will be non-nil
//If the backend malfunctions, user will be nil and error will be non-nil.
func (a *LDAPAuth) AdminLogin(username, password string) (user *User, err error) {
	ok, err := auth.Login(username, password, a.adminGroup, a.config)
	if !ok || err != nil {
		return nil, err
	}

	return &User{Username: username, Admin: true}, nil
}
