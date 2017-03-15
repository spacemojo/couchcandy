package couchcandy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// CouchCandy : Struct that provides all CouchDB's API has to offer.
type CouchCandy struct {
	LclSession *Session
	GetHandler func(string) (*http.Response, error)
}

// NewCouchCandy : Returns a new CouchCandy struct initialised with the provided values.
func NewCouchCandy(session *Session) *CouchCandy {
	cc := &CouchCandy{
		LclSession: session,
		GetHandler: http.Get,
	}
	return cc
}

// CreateDatabaseURL : Creates the right URL to point to the passed database.
func CreateDatabaseURL(session *Session, db string) string {
	var buffer bytes.Buffer
	buffer.WriteString(session.Host)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(session.Port))
	buffer.WriteString("/")
	buffer.WriteString(db)
	return buffer.String()
}

// GetDatabaseInfo returns basic information about the passed database.
func (c *CouchCandy) GetDatabaseInfo(db string) (*DatabaseInfo, error) {

	url := CreateDatabaseURL(c.LclSession, db)
	res, geterr := c.GetHandler(url)
	if geterr != nil {
		return nil, geterr
	}

	page, geterr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if geterr != nil {
		return nil, geterr
	}

	dbInfo := &DatabaseInfo{}
	json.Unmarshal(page, dbInfo)
	return dbInfo, nil

}
