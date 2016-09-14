package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

//DB represents a database
type DB interface {
	//Submit commits an Entry to the database, and returns an error if one occurred
	Submit(e *Entry) error

	//Check returns whether or not the given employeeID is already in the database.
	//Check returns an error if one occurred.
	Check(employeeID string) (bool, error)

	//List returns all entries in the database
	List() ([]*Entry, error)
}

//SQLDB is a DB backed by a SQL database
type SQLDB struct {
	db     *sql.DB
	driver string
	dsn    string
}

//Submit commits an Entry to the database, and returns an error if one occurred
func (db *SQLDB) Submit(e *Entry) error {
	j, err := json.Marshal(e.Headers)
	if err != nil {
		return err
	}
	_, err = db.db.Exec("INSERT INTO signers(employeeID, username, firstname, lastname, campus, headers, time) VALUES(?, ?, ?, ?, ?, ?);",
		e.EmployeeID,
		e.Username,
		e.FirstName,
		e.LastName,
		e.Campus,
		j,
		e.Time,
	)
	return err
}

//Check returns whether or not the given employeeID is already in the database
//Check returns an error if one occurred.
func (db *SQLDB) Check(employeeID string) (bool, error) {
	row := db.db.QueryRow("SELECT employee_id FROM signers WHERE employee_id=?;", employeeID)

	s := new(string)
	err := row.Scan(s)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

//List returns all entries in the database
func (db *SQLDB) List() (list []*Entry, err error) {
	rows, err := db.db.Query("SELECT employee_id, username, firstname, lastname, campus, headers, time FROM signers;")
	if err != nil {
		return nil, err
	}

	defer func() {
		e := rows.Close()
		if err == nil {
			err = e
		}
	}()

	var entries []*Entry
	for rows.Next() {
		e := &Entry{}
		var j []byte

		err = rows.Scan(&(e.EmployeeID), &(e.Username), &(e.FirstName), &(e.LastName), &(e.Campus), &j, &(e.Time))
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(j, &(e.Headers))
		if err != nil {
			return nil, err
		}

		entries = append(entries, e)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

//NewSQLDB creates a new SQLDB with the given driver and dsn as used by database/sql's Open
func NewSQLDB(driver, dsn string) (*SQLDB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &SQLDB{
		db:     db,
		driver: driver,
		dsn:    dsn,
	}, nil
}

//Entry represents a database entry
type Entry struct {
	EmployeeID string
	Username   string
	FirstName  string
	LastName   string
	Campus     string
	Headers    http.Header
	Time       time.Time
}

//Validate makes sure the information in e is valid
func (e *Entry) Validate() error {
	if len(e.Username) > 255 {
		return errors.New("Username > 255")
	}
	if len(e.FirstName) > 255 {
		return errors.New("FirstName > 255")
	}
	if len(e.LastName) > 255 {
		return errors.New("LastName > 255")
	}
	switch e.Campus {
	case "Early Childhood", "Primary", "Elementary", "Intermediate", "Middle", "High", "Central Office", "Transportation", "Maintenance":
	default:
		return fmt.Errorf("Undefined Campus: %s", e.Campus)
	}
	return nil
}

//NewEntry creates a new Entry with the given information
func NewEntry(u *User, s *SubmitRequest, h http.Header) *Entry {
	return &Entry{
		EmployeeID: u.EmployeeID,
		Username:   u.Username,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Campus:     s.Campus,
		Headers:    h,
		Time:       time.Now(),
	}
}
