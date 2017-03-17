package couchcandy

import (
	"fmt"
	"strconv"
)

// CreateBaseURL : Creates the base of all urls all the way up to the port number.
func CreateBaseURL(session Session) string {
	return fmt.Sprintf("http://%s:%s@%s:%s", session.Username, session.Password, session.Host[7:], strconv.Itoa(session.Port))
}

// CreateDatabaseURL : Creates the right URL to point to the passed database.
func CreateDatabaseURL(session Session) string {
	return fmt.Sprintf("%s/%s", CreateBaseURL(session), session.Database)
}

// CreateDocumentURL : Creates the right URL to fetch the passed document in the passed database.
func CreateDocumentURL(session Session, id string) string {
	return fmt.Sprintf("%s/%s", CreateDatabaseURL(session), id)
}

// CreateAllDatabasesURL : Creates the URL that allows to fetch all the
// database names in CouchDB
func CreateAllDatabasesURL(session Session) string {
	return fmt.Sprintf("%s/_all_dbs", CreateBaseURL(session))
}
