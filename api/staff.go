package api

import (
	"database/sql"
	"fmt"
	"strings"
)

const staffQuery = `
SELECT
    name."ALTERNATE-ID" AS EmployeeID,
    name."FIRST-NAME" AS FirstName,
    name."LAST-NAME" AS LastName,
    empcode."HAAETY-DESC" AS EmployeeTypeDesc,
    bldcode."HAABLD-DESC" AS BuildingDesc
FROM PUB.NAME AS name
    INNER JOIN PUB."HAAPRO-PROFILE" AS profile
ON 
    name.NALPHAKEY = profile.nalphakey
INNER JOIN PUB."HAAETY-EMP-TYPES" AS empcode ON
    empcode."HAAETY-EMP-TYPE-CODE" = profile."HAAETY-EMP-TYPE-CODE"
INNER JOIN PUB."HAABLD-BLD-CODES" as bldcode ON
    bldcode."HAABLD-BLD-CODE" = profile."HAABLD-BLD-CODE"
WHERE empcode."HAAETY-EMP-TYPE-CODE" <> 'TERM'
WITH (READPAST NOWAIT)
`

//StaffMember represents information about a staff member
type StaffMember struct {
	EmployeeID string
	FirstName  string
	LastName   string
	Type       string
	Location   string
}

//StaffDB represents a Staff database
type StaffDB interface {
	//List returns all entries in the database
	List() ([]*StaffMember, error)
}

//SkywardDB is a DB backed by a Skyward (Progress) Database
type SkywardDB struct {
	db     *sql.DB
	driver string
	dsn    string
	skips  map[string]struct{}
}

//NewSkywardDB creates a new SkywardDB with the given driver and dsn as used by database/sql's Open.
//If excludedIDs is not nil, then any StaffMembers with an EmployeeID in excludedIDs will not be returned by List.
func NewSkywardDB(driver, dsn string, excludedIDs []string) (*SkywardDB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	skips := make(map[string]struct{})
	for _, skip := range excludedIDs {
		skips[skip] = struct{}{}
	}
	return &SkywardDB{
		db:     db,
		driver: driver,
		dsn:    dsn,
		skips:  skips,
	}, nil
}

//List returns all entries in the database
func (db *SkywardDB) List() (list []*StaffMember, err error) {
	rows, err := db.db.Query(staffQuery)
	if err != nil {
		return nil, err
	}

	defer func() {
		e := rows.Close()
		if err == nil {
			err = e
		}
	}()

	var staff []*StaffMember
	for rows.Next() {
		s := &StaffMember{}
		var id int64

		err = rows.Scan(&id, &(s.FirstName), &(s.LastName), &(s.Type), &(s.Location))
		if err != nil {
			return nil, err
		}

		s.FirstName = strings.Title(strings.TrimSpace(s.FirstName))
		s.LastName = strings.Title(strings.TrimSpace(s.LastName))
		s.Type = strings.Title(strings.TrimSpace(s.Type))
		s.Location = strings.Title(strings.TrimSpace(s.Location))

		s.EmployeeID = fmt.Sprintf("f%d", id)

		if _, ok := db.skips[s.EmployeeID]; !ok {
			staff = append(staff, s)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return staff, nil
}
