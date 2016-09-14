package api

import "time"

//AuthRequest is a client->server request for authentication
type AuthRequest struct {
	User   string
	Passwd string
}

//AuthResponse is a server->client response about authentication
type AuthResponse struct {
	SessionID string
	Completed bool //if username exists in database
}

//SubmitRequest is a client->server request for submitting form information
type SubmitRequest struct {
	Campus string
	Agree  bool
}

//SubmitResponse is a server->client response about confirming a submission
type SubmitResponse struct {
	Status bool
}

//ErrorResponse is a server-client response indicating some kind of error
type ErrorResponse struct {
	Code  int
	Error string
}

//Record represents a staff signing record
type Record struct {
	FirstName    string
	LastName     string
	EmployeeType string
	Location     string
	SignTime     time.Time
}

//ListResponse is a server->client response with a list of signing records
type ListResponse struct {
	List []*Record
}
