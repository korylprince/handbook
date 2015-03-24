package api

import (
	"net/http"
)

//Context represents a group of services
type Context struct {
	Auth         Auth
	DB           DB
	SessionStore SessionStore
}

type contextHandler struct {
	HandleFunc func(*Context, http.ResponseWriter, *http.Request)
	Context    *Context
}

func (c contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.HandleFunc(c.Context, w, r)
}

//AuthHandler returns an Authentcation http.Handler with the given context
func AuthHandler(c *Context) http.Handler {
	return contextHandler{HandleFunc: authNormalHandler, Context: c}
}

//AuthAdminHandler returns an Admin Authentcation http.Handler with the given context
func AuthAdminHandler(c *Context) http.Handler {
	return contextHandler{HandleFunc: authAdminHandler, Context: c}
}

//SubmitHandler returns a Submission http.Handler with the given context
func SubmitHandler(c *Context) http.Handler {
	return contextHandler{HandleFunc: submitHandler, Context: c}
}

//ListHandler returns a dump of the given context's DB
func ListHandler(c *Context) http.Handler {
	return contextHandler{HandleFunc: listHandler, Context: c}
}
