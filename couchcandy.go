package couchcandy

import (
	"encoding/json"
	"io/ioutil"
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

	res, puterr := c.PutHandler(url, body)
	if puterr != nil {
		return nil, puterr
	}

	page, puterr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if puterr != nil {
		return nil, puterr
	}

	return page, nil

}

func (c *CouchCandy) readFromGet(url string) ([]byte, error) {

	res, geterr := c.GetHandler(url)
	if geterr != nil {
		return nil, geterr
	}

	page, geterr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if geterr != nil {
		return nil, geterr
	}

	return page, nil

}
