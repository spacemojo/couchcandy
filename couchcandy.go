package couchcandy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// DatabaseInfo returns basic information about the database in session.
func (c *CouchCandy) DatabaseInfo() (*DatabaseInfo, error) {

	url := createDatabaseURL(c.Session)
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return nil, err
	}

	dbInfo := &DatabaseInfo{}
	unmarshallError := json.Unmarshal(page, dbInfo)
	return dbInfo, unmarshallError

}

// Document Returns the specified document.
func (c *CouchCandy) Document(id string, v interface{}, options Options) error {

	url := createDocumentURLWithOptions(c.Session, id, options)
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return err
	}

	unmarshallError := json.Unmarshal(page, v)
	return unmarshallError

}

// Add Adds a document in the database but the system will generate
// an id. Look at PutDocumentWithID for setting an id for the document explicitly.
func (c *CouchCandy) Add(document interface{}) (*OperationResponse, error) {

	url := createDatabaseURL(c.Session)
	bodyStr, marshallError := safeMarshall(document)
	if marshallError != nil {
		return nil, marshallError
	}

	page, err := readFromWithBody(url, bodyStr, c.PostHandler)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// Update Updates a document in the database. Note that _id and _rev
// fields are required in the passed document.
func (c *CouchCandy) Update(document interface{}) (*OperationResponse, error) {

	bodyStr, marshallError := safeMarshall(document)
	if marshallError != nil {
		return nil, marshallError
	}

	url := createPutDocumentURL(c.Session, bodyStr)

	page, err := readFromWithBody(url, bodyStr, c.PutHandler)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// AddWithID Inserts a document in the database with the specified id
func (c *CouchCandy) AddWithID(id string, document interface{}) (*OperationResponse, error) {

	url := fmt.Sprintf("%s/%s", createDatabaseURL(c.Session), id)

	bodyStr, marshallError := safeMarshall(document)
	if marshallError != nil {
		return nil, marshallError
	}

	page, err := readFromWithBody(url, bodyStr, c.PutHandler)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// AddAttachment adds the provided attachment to the specified
// document and revision.
func (c *CouchCandy) AddAttachment(id, rev, name, contentType string, file []byte) (*OperationResponse, error) {

	url := fmt.Sprintf("%s/%s/%s?rev=%s", createDatabaseURL(c.Session), id, name, rev)
	fmt.Printf("Attachment url : %s\n", url)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(file))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", contentType)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	page, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// Documents : Returns all documents in the database based on the passed parameters.
func (c *CouchCandy) Documents(options Options) (*AllDocuments, error) {

	url := fmt.Sprintf("%s/_all_docs%s", createDatabaseURL(c.Session), toQueryString(options))
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return nil, err
	}

	return toAllDocuments(page)

}

// DocumentsByKeys Fetches all the documents corresponding to the passed keys array.
func (c *CouchCandy) DocumentsByKeys(keys []string, options Options) (*AllDocuments, error) {

	url := fmt.Sprintf("%s/_all_docs%s", createDatabaseURL(c.Session), toQueryString(options))

	body, _ := json.Marshal(&AllDocumentsKeys{
		Keys: keys,
	})

	page, err := readFromWithBody(url, string(body), c.PostHandler)
	if err != nil {
		return nil, err
	}

	return toAllDocuments(page)

}

// AddDatabase : Creates a database in CouchDB
func (c *CouchCandy) AddDatabase(name string) (*OperationResponse, error) {

	c.Session.Database = name
	url := createDatabaseURL(c.Session)

	page, err := readFromWithBody(url, "", c.PutHandler)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// DeleteDatabase : Deletes the passed database from the system.
func (c *CouchCandy) DeleteDatabase(name string) (*OperationResponse, error) {

	c.Session.Database = name
	url := createDatabaseURL(c.Session)
	page, err := readFrom(url, c.DeleteHandler)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// Delete Deletes the passed document with revision from the database
func (c *CouchCandy) Delete(id string, revision string) (*OperationResponse, error) {

	url := fmt.Sprintf("%s?rev=%s", createDocumentURL(c.Session, id), revision)
	page, err := readFrom(url, c.DeleteHandler)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// AllDatabases : Returns all the database names in the system.
func (c *CouchCandy) AllDatabases() ([]string, error) {

	url := createAllDatabasesURL(c.Session)
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return nil, err
	}

	var dbs []string
	unmarshallError := json.Unmarshal(page, &dbs)
	return dbs, unmarshallError

}

// ChangeNotifications : Return the current change notifications.
func (c *CouchCandy) ChangeNotifications(options Options) (*Changes, error) {

	url := fmt.Sprintf("%s/_changes?style=%s", createDatabaseURL(c.Session), options.Style)
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return nil, err
	}

	changes := &Changes{}
	unmarshallError := json.Unmarshal(page, changes)
	return changes, unmarshallError

}

// View : Calls the passed view with provided options
func (c *CouchCandy) View(ddoc, view string, options Options) (*ViewResponse, error) {

	url := fmt.Sprintf("%s/_design/%s/_view/%s%s", createDatabaseURL(c.Session), ddoc, view, toQueryString(options))
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return nil, err
	}
	return toViewResponse(page)

}

// ViewWithList calls the passed view with list and options
func (c *CouchCandy) ViewWithList(ddoc, list, view string, options Options) (*ViewResponse, error) {

	url := fmt.Sprintf("%s/_design/%s/_list/%s/%s/%s%s", createDatabaseURL(c.Session), ddoc, list, ddoc, view, toQueryString(options))
	fmt.Printf("CouchCandy.CallView(%s)\n", url)
	page, err := readFrom(url, c.GetHandler)
	if err != nil {
		return nil, err
	}
	return toViewResponse(page)

}

func readFromWithBody(url, body string, handler func(str string, bd string) (*http.Response, error)) ([]byte, error) {

	res, err := handler(url, body)
	if err != nil {
		return nil, err
	}

	page, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
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
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return page, nil

}

func defaultPostHandler(url, body string) (*http.Response, error) {
	return defaultHandlerWithBody(http.MethodPost, url, body, &http.Client{})
}

func defaultPutHandler(url, body string) (*http.Response, error) {
	return defaultHandlerWithBody(http.MethodPut, url, body, &http.Client{})
}

func defaultHandlerWithBody(method, url, body string, client CandyHTTPClient) (*http.Response, error) {

	bodyJSON := strings.NewReader(body)
	request, requestError := http.NewRequest(method, url, bodyJSON)
	if requestError != nil {
		return nil, requestError
	}

	request.Header.Add("Content-Type", "application/json")
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
