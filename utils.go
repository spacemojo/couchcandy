package couchcandy

import (
	"bytes"
	"strconv"
)

// CreateBaseURL : Creates the base of all urls all the way up to the port number.
func CreateBaseURL(session *Session) string {
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(session.Username)
	buffer.WriteString(":")
	buffer.WriteString(session.Password)
	buffer.WriteString("@")
	buffer.WriteString(session.Host[7:])
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(session.Port))
	return buffer.String()
}

// CreateDatabaseURL : Creates the right URL to point to the passed database.
func CreateDatabaseURL(session *Session) string {
	var buffer bytes.Buffer
	buffer.WriteString(CreateBaseURL(session))
	buffer.WriteString("/")
	buffer.WriteString(session.Database)
	return buffer.String()
}

// CreateDocumentURL : Creates the right URL to fetch the passed document in the passed database.
func CreateDocumentURL(session *Session, id string) string {
	var buffer bytes.Buffer
	buffer.WriteString(CreateDatabaseURL(session))
	buffer.WriteString("/")
	buffer.WriteString(id)
	return buffer.String()
}

// CreateAllDatabasesURL : Creates the URL that allows to fetch all the
// database names in CouchDB
func CreateAllDatabasesURL(session *Session) string {
	var buffer bytes.Buffer
	buffer.WriteString(CreateBaseURL(session))
	buffer.WriteString("/_all_dbs")
	return buffer.String()
}
