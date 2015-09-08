package api

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

//ListResponse is a server->client response with a list of entries
type ListResponse struct {
	List []*Entry
}

//MissingListResponse is a server->client response with a list of LDAP Users who have not signed
type MissingListResponse struct {
	List []*LDAPUser
}
