package couchcandy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GetDatabaseInfo returns basic information about the database in session.
func (c *CouchCandy) GetDatabaseInfo() (*DatabaseInfo, error) {

	url := CreateDatabaseURL(c.LclSession)
	page, err := c.readFromGet(url)
	if err != nil {
		return nil, err
	}

	dbInfo := &DatabaseInfo{}
	unmarshallError := json.Unmarshal(page, dbInfo)
	return dbInfo, unmarshallError

}

// GetDocument : Returns the specified document in the passed database.
func (c *CouchCandy) GetDocument(id string, v interface{}) error {

	url := CreateDocumentURL(c.LclSession, id)
	page, err := c.readFromGet(url)
	if err != nil {
		return err
	}

	unmarshallError := json.Unmarshal(page, v)
	return unmarshallError

}

// PutDatabase : Creates a database in CouchDB
func (c *CouchCandy) PutDatabase(name string) (*OperationResponse, error) {

	c.LclSession.Database = name
	url := CreateDatabaseURL(c.LclSession)

	page, err := c.readFromPut(url, "")
	if err != nil {
		return nil, err
	}

	response := &OperationResponse{}
	unmarshallError := json.Unmarshal(page, response)
	return response, unmarshallError

}

// DeleteDatabase : Deletes the passed database from the system.
func (c *CouchCandy) DeleteDatabase(name string) (*OperationResponse, error) {

	c.LclSession.Database = name
	url := CreateDatabaseURL(c.LclSession)
	page, err := c.readFromDelete(url)
	if err != nil {
		return nil, err
	}

	response := &OperationResponse{}
	unmarshallError := json.Unmarshal(page, response)
	return response, unmarshallError

}

// GetAllDatabases : Returns all the database names in the system.
func (c *CouchCandy) GetAllDatabases() ([]string, error) {

	url := CreateAllDatabasesURL(c.LclSession)
	page, err := c.readFromGet(url)
	if err != nil {
		return nil, err
	}

	var dbs []string
	unmarshallError := json.Unmarshal(page, &dbs)
	return dbs, unmarshallError

}

func (c *CouchCandy) readFromPut(url string, body string) ([]byte, error) {
	return readFromWithBody(url, body, c.PutHandler)
}

func (c *CouchCandy) readFromDelete(url string) ([]byte, error) {
	return readFrom(url, c.DeleteHandler)
}

func (c *CouchCandy) readFromGet(url string) ([]byte, error) {
	return readFrom(url, c.GetHandler)
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
