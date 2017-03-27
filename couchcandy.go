package couchcandy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetDatabaseInfo returns basic information about the database in session.
func (c *CouchCandy) GetDatabaseInfo() (*DatabaseInfo, error) {

	url := createDatabaseURL(c.LclSession)
	page, err := c.readFromGet(url)
	if err != nil {
		return nil, err
	}

	dbInfo := &DatabaseInfo{}
	unmarshallError := json.Unmarshal(page, dbInfo)
	return dbInfo, unmarshallError

}

// GetDocument Returns the specified document.
func (c *CouchCandy) GetDocument(id string, v interface{}, options Options) error {

	url := createDocumentURLWithOptions(c.LclSession, id, options)
	page, err := c.readFromGet(url)
	if err != nil {
		return err
	}

	unmarshallError := json.Unmarshal(page, v)
	return unmarshallError

}

// PostDocument Adds a document in the database but the system will generate
// an id. Look at PutDocumentWithID for setting an id for the document explicitly.
func (c *CouchCandy) PostDocument(document interface{}) (*OperationResponse, error) {

	url := createDatabaseURL(c.LclSession)

	body, marshallError := json.Marshal(document)
	if marshallError != nil {
		return nil, marshallError
	}

	page, err := c.readFromPost(url, string(body))
	if err != nil {
		return nil, err
	}

	response := &OperationResponse{}
	unmarshallError := json.Unmarshal(page, response)
	return response, unmarshallError

}

// PutDocument Updates a document in the database. Note that _id and _rev
// fields are required in the passed document.
func (c *CouchCandy) PutDocument(document interface{}) (*OperationResponse, error) {

	url := createDatabaseURL(c.LclSession)
	body, marshallError := json.Marshal(document)
	if marshallError != nil {
		return nil, marshallError
	}

	page, err := c.readFromPut(url, string(body))
	if err != nil {
		return nil, err
	}

	response := &OperationResponse{}
	unmarshallError := json.Unmarshal(page, response)
	return response, unmarshallError

}

// PutDocumentWithID Inserts a document in the database with the specified id
func (c *CouchCandy) PutDocumentWithID(id string, document interface{}) (*OperationResponse, error) {

	url := fmt.Sprintf("%s/%s", createDatabaseURL(c.LclSession), id)

	body, marshallError := json.Marshal(document)
	if marshallError != nil {
		return nil, marshallError
	}

	page, err := c.readFromPut(url, string(body))
	if err != nil {
		return nil, err
	}

	response := &OperationResponse{}
	unmarshallError := json.Unmarshal(page, response)
	return response, unmarshallError

}

// GetAllDocuments : Returns all documents in the database based on the passed parameters.
func (c *CouchCandy) GetAllDocuments(options Options) (*AllDocuments, error) {

	url := fmt.Sprintf("%s/_all_docs?descending=%v&limit=%v&include_docs=%v", createDatabaseURL(c.LclSession), options.Descending, options.Limit, options.IncludeDocs)
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return nil, err
	}

	allDocuments := &AllDocuments{}
	unmarshallError := json.Unmarshal(page, allDocuments)
	return allDocuments, unmarshallError

}

// PutDatabase : Creates a database in CouchDB
func (c *CouchCandy) PutDatabase(name string) (*OperationResponse, error) {

	c.LclSession.Database = name
	url := createDatabaseURL(c.LclSession)

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
	url := createDatabaseURL(c.LclSession)
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

	url := createAllDatabasesURL(c.LclSession)
	page, err := c.readFromGet(url)
	if err != nil {
		return nil, err
	}

	var dbs []string
	unmarshallError := json.Unmarshal(page, &dbs)
	return dbs, unmarshallError

}

// GetChangeNotifications : Return the current change notifications.
func (c *CouchCandy) GetChangeNotifications(options Options) (*Changes, error) {

	url := fmt.Sprintf("%s/_changes?style=%s", createDatabaseURL(c.LclSession), options.Style)
	page, err := c.readFromGet(url)
	if err != nil {
		return nil, err
	}

	changes := &Changes{}
	unmarshallError := json.Unmarshal(page, changes)
	return changes, unmarshallError

}

func (c *CouchCandy) readFromPut(url string, body string) ([]byte, error) {
	return readFromWithBody(url, body, c.PutHandler)
}

func (c *CouchCandy) readFromPost(url string, body string) ([]byte, error) {
	return readFromWithBody(url, body, c.PostHandler)
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
