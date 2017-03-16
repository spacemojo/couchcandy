package couchcandy

import (
	"bytes"
	"io/ioutil"
	"net/http"
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

func (c *CouchCandy) readFromPut(url string, body string) ([]byte, error) {
	return readFromWithBody(url, body, c.PutHandler)
}

func readFromWithBody(url string, body string, handler func(str string, bd string) (*http.Response, error)) ([]byte, error) {

	res, err := handler(url, body)
	if err != nil {
		return nil, err
	}

	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return page, nil

}

func readFrom(url string, handler func(str string) (*http.Response, error)) ([]byte, error) {

	res, err := handler(url)
	if err != nil {
		return nil, err
	}

	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return page, nil

}

func (c *CouchCandy) readFromDelete(url string) ([]byte, error) {
	return readFrom(url, c.DeleteHandler)
}

func (c *CouchCandy) readFromGet(url string) ([]byte, error) {
	return readFrom(url, c.GetHandler)
}
