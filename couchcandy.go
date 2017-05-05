package couchcandy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetDatabaseInfo returns basic information about the database in session.
func (c *CouchCandy) GetDatabaseInfo() (*DatabaseInfo, error) {

	url := createDatabaseURL(c.LclSession)
	page, err := readFrom(url, c.GetHandler)
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
	page, err := readFrom(url, c.GetHandler)
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
	bodyStr, marshallError := safeMarshall(document)
	if marshallError != nil {
		return nil, marshallError
	}

	page, err := readFromWithBody(url, bodyStr, c.PostHandler)
	if err != nil {
		return nil, err
	}

	return produceOperationResponse(page)

}

// PutDocument Updates a document in the database. Note that _id and _rev
// fields are required in the passed document.
func (c *CouchCandy) PutDocument(document interface{}) (*OperationResponse, error) {

	url := createDatabaseURL(c.LclSession)
	bodyStr, marshallError := safeMarshall(document)
	if marshallError != nil {
		return nil, marshallError
	}

	page, err := readFromWithBody(url, bodyStr, c.PutHandler)
	if err != nil {
		return nil, err
	}

	return produceOperationResponse(page)

}

// PutDocumentWithID Inserts a document in the database with the specified id
func (c *CouchCandy) PutDocumentWithID(id string, document interface{}) (*OperationResponse, error) {

	url := fmt.Sprintf("%s/%s", createDatabaseURL(c.LclSession), id)

	bodyStr, marshallError := safeMarshall(document)
	if marshallError != nil {
		return nil, marshallError
	}

	page, err := readFromWithBody(url, bodyStr, c.PutHandler)
	if err != nil {
		return nil, err
	}

	return produceOperationResponse(page)

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

	page, err := readFromWithBody(url, "", c.PutHandler)
	if err != nil {
		return nil, err
	}

	return produceOperationResponse(page)

}

// DeleteDatabase : Deletes the passed database from the system.
func (c *CouchCandy) DeleteDatabase(name string) (*OperationResponse, error) {

	c.LclSession.Database = name
	url := createDatabaseURL(c.LclSession)
	page, err := readFrom(url, c.DeleteHandler)
	if err != nil {
		return nil, err
	}

	return produceOperationResponse(page)

}

// DeleteDocument Deletes the passed document with revision from the database
func (c *CouchCandy) DeleteDocument(id string, revision string) (*OperationResponse, error) {

	url := fmt.Sprintf("%s?rev=%s", createDocumentURL(c.LclSession, id), revision)
	page, err := readFrom(url, c.DeleteHandler)
	if err != nil {
		return nil, err
	}

	return produceOperationResponse(page)

}

// GetAllDatabases : Returns all the database names in the system.
func (c *CouchCandy) GetAllDatabases() ([]string, error) {

	url := createAllDatabasesURL(c.LclSession)
	page, err := readFrom(url, c.GetHandler)
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
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return nil, err
	}

	changes := &Changes{}
	unmarshallError := json.Unmarshal(page, changes)
	return changes, unmarshallError

}

func produceOperationResponse(page []byte) (*OperationResponse, error) {
	response := &OperationResponse{}
	unmarshallError := json.Unmarshal(page, response)
	return response, unmarshallError
}

// this is a violent hack to set the Revisions field to nil so that it does no get marshalled initially.
func safeMarshall(document interface{}) (string, error) {
	body, err := json.Marshal(document)
	if err != nil {
		return "", err
	}
	bodyStr := strings.Replace(string(body), "\"_revisions\":{\"start\":0,\"ids\":null},", "", -1)
	return bodyStr, nil
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

func defaultPostHandler(url string, body string) (*http.Response, error) {
	return defaultHandlerWithBody(http.MethodPost, url, body, &http.Client{})
}

func defaultPutHandler(url string, body string) (*http.Response, error) {
	return defaultHandlerWithBody(http.MethodPut, url, body, &http.Client{})
}

func defaultHandlerWithBody(method, url, body string, client CandyHTTPClient) (*http.Response, error) {

	bodyJson := strings.NewReader(body)
	fmt.Printf("JSON BODY : %s\n", bodyJson)

	request, requestError := http.NewRequest(method, url, bodyJson)
	request.Header.Add("Content-Type", "application/json")
	if requestError != nil {
		return nil, requestError
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func defaultGetHandler(url string) (*http.Response, error) {
	return defaultHandler(http.MethodGet, url, &http.Client{})
}

func defaultDeleteHandler(url string) (*http.Response, error) {
	return defaultHandler(http.MethodDelete, url, &http.Client{})
}

func defaultHandler(method, url string, client CandyHTTPClient) (*http.Response, error) {

	request, requestError := http.NewRequest(method, url, nil)
	if requestError != nil {
		return nil, requestError
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil

}
