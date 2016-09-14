package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

//handleError returns a json response for the given code and logs the error
func handleError(w http.ResponseWriter, code int, err error) {
	log.Println(err)
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	encErr := e.Encode(ErrorResponse{Code: code, Error: http.StatusText(code)})
	if encErr != nil {
		panic(encErr)
	}
}

//NotFoundHandler returns a json 401 response
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	handleError(w, http.StatusNotFound, errors.New("handler not found"))
}

//authHandler will return a sessionID if the credentials are valid
//or an HTTP 401 Error if not.
//The admin flag specifies whether to use Login or AdminLogin functions
func authHandler(admin bool, c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var aReq AuthRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&aReq)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("Error decoding json: %v", err))
		return
	}

	var user *User
	if admin {
		user, err = c.Auth.AdminLogin(aReq.User, aReq.Passwd)
	} else {
		user, err = c.Auth.Login(aReq.User, aReq.Passwd)
	}
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Error authenticating: %v", err))
		return
	}
	if user != nil {
		sessionID, err := c.SessionStore.Create(user)
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error creating session key: %v", err))
			return
		}

		completed, err := c.DB.Check(aReq.User)
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error checking database for username: %v", err))
			return
		}

		e := json.NewEncoder(w)
		err = e.Encode(AuthResponse{SessionID: sessionID, Completed: completed})
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error encoding json: %v", err))
		}
		return
	}
	handleError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
}

//authNormalHandler will return a sessionID if the credentials are a valid login
//or an HTTP 401 Error if not.
func authNormalHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	authHandler(false, c, w, r)
}

//authAdminHandler will return a sessionID if the credentials are a valid admin login
//or an HTTP 401 Error if not.
func authAdminHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	authHandler(true, c, w, r)
}

//submitHandler will submit information to the database if the sessionID is valid
//or an HTTP 401 Error if not.
func submitHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	key := r.Header.Get("X-Session-Key")
	if key == "" {
		handleError(w, http.StatusBadRequest, errors.New("X-Session-Key header empty"))
		return
	}

	var sReq SubmitRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&sReq)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("Error decoding json: %v", err))
		return
	}

	if !sReq.Agree {
		handleError(w, http.StatusBadRequest, errors.New("SubmitRequest.Agree was false"))
		return
	}

	sess, err := c.SessionStore.Check(key)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Error checking session key: %v", err))
		return
	}
	if sess != nil {
		entry := NewEntry(sess.User, &sReq, r.Header)

		err = entry.Validate()
		if err != nil {
			handleError(w, http.StatusBadRequest, fmt.Errorf("Error validating entry: %v", err))
			return
		}

		err = c.DB.Submit(NewEntry(sess.User, &sReq, r.Header))
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error submitting entry to database: %v", err))
			return
		}

		e := json.NewEncoder(w)
		err = e.Encode(SubmitResponse{Status: true})
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error encoding json: %v", err))
		}
		return
	}
	handleError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
}

//listHandler will return a dump of the database if the sessionID is valid
//or an HTTP 401 Error if not.
func listHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	key := r.Header.Get("X-Session-Key")
	if key == "" {
		handleError(w, http.StatusBadRequest, errors.New("X-Session-Key header empty"))
		return
	}

	sess, err := c.SessionStore.Check(key)
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("Error checking session key: %v", err))
		return
	}
	if sess != nil && sess.User != nil && sess.User.Admin {

		list, err := c.DB.List()
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error getting list from database: %v", err))
			return
		}

		e := json.NewEncoder(w)
		err = e.Encode(ListResponse{List: list})
		if err != nil {
			handleError(w, http.StatusInternalServerError, fmt.Errorf("Error encoding json: %v", err))
		}
		return
	}
	handleError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
}
