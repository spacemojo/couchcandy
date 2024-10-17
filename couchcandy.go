package couchcandy

import (
	"encoding/json"
	"fmt"
)

// DatabaseInfo returns basic information about the database in session.
func (c *CouchCandy) DatabaseInfo() (*DatabaseInfo, error) {

	url := createDatabaseURL(c.Session)
	page, err := readJSON(url, c.Get)
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
	page, err := readJSON(url, c.Get)
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

	page, err := readJSONWithBody(url, bodyStr, c.PostJSON)
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

	page, err := readJSONWithBody(url, bodyStr, c.PutJSON)
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

	page, err := readJSONWithBody(url, bodyStr, c.PutJSON)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// AddAttachment adds the provided attachment to the specified
// document and revision.
func (c *CouchCandy) AddAttachment(id, rev, name, contentType string, file []byte) (*OperationResponse, error) {

	url := fmt.Sprintf("%s/%s/%s?rev=%s", createDatabaseURL(c.Session), id, name, rev)

	page, err := readBytesWithBody(url, contentType, file, c.PutBytes)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// DeleteAttachment delets the named attachment on the document corresponding
// to the id-rev pair.
func (c *CouchCandy) DeleteAttachment(id, rev, name string) (*OperationResponse, error) {

	// DELETE /db/doc/attachmentname?rev=...
	url := fmt.Sprintf("%s/%s/%s?rev=%s", createDatabaseURL(c.Session), id, name, rev)

	page, err := readJSON(url, c.Delete)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// Documents : Returns all documents in the database based on the passed parameters.
func (c *CouchCandy) Documents(options Options) (*AllDocuments, error) {

	url := fmt.Sprintf("%s/_all_docs%s", createDatabaseURL(c.Session), toQueryString(options))
	page, err := readJSON(url, c.Get)
	if err != nil {
		return nil, err
	}

	return toAllDocuments(page)

}

func (c *CouchCandy) DesignDocs() (*DesignDocs, error) {

	allDocuments, err := c.Documents(Options{
		StartKey: fmt.Sprintf("\"%s\"", "_design"),
		EndKey:   fmt.Sprintf("\"%s\"", "_design0"),
	})

	if err != nil {
		return nil, err
	}

	designDocs := NewDesignDocs()

	for _, doc := range allDocuments.Rows {

		partial := &partialDesignDoc{}

		err = json.Unmarshal(doc.Doc, partial)
		if err != nil {
			return nil, err
		}

		if partial.Language == "query" {

			index := &IndexDesignDoc{}
			err = json.Unmarshal(doc.Doc, index)
			if err != nil {
				return nil, err
			}
			designDocs.Index = append(designDocs.Index, index)

		} else if partial.Language == "javascript" {

			mapReduce := &MapReduceDesignDoc{}
			err = json.Unmarshal(doc.Doc, mapReduce)
			if err != nil {
				return nil, err
			}
			designDocs.MapReduce = append(designDocs.MapReduce, mapReduce)

		}

	}

	return designDocs, nil

}

// DocumentsByKeys Fetches all the documents corresponding to the passed keys array.
func (c *CouchCandy) DocumentsByKeys(keys []string, options Options) (*AllDocuments, error) {

	url := fmt.Sprintf("%s/_all_docs%s", createDatabaseURL(c.Session), toQueryString(options))

	body, _ := json.Marshal(&AllDocumentsKeys{
		Keys: keys,
	})

	page, err := readJSONWithBody(url, string(body), c.PostJSON)
	if err != nil {
		return nil, err
	}

	return toAllDocuments(page)

}

// AddDatabase : Creates a database in CouchDB
func (c *CouchCandy) AddDatabase(name string) (*OperationResponse, error) {

	c.Session.Database = name
	url := createDatabaseURL(c.Session)

	page, err := readJSONWithBody(url, "", c.PutJSON)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// DeleteDatabase : Deletes the passed database from the system.
func (c *CouchCandy) DeleteDatabase(name string) (*OperationResponse, error) {

	c.Session.Database = name
	url := createDatabaseURL(c.Session)
	page, err := readJSON(url, c.Delete)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// DeleteDocument Deletes the passed document with revision from the database
func (c *CouchCandy) DeleteDocument(id string, revision string) (*OperationResponse, error) {

	url := fmt.Sprintf("%s?rev=%s", createDocumentURL(c.Session, id), revision)
	page, err := readJSON(url, c.Delete)
	if err != nil {
		return nil, err
	}

	return toOperationResponse(page)

}

// AllDatabases : Returns all the database names in the system.
func (c *CouchCandy) AllDatabases() ([]string, error) {

	url := createAllDatabasesURL(c.Session)
	page, err := readJSON(url, c.Get)
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
	page, err := readJSON(url, c.Get)
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
	page, err := readJSON(url, c.Get)
	if err != nil {
		return nil, err
	}
	return toViewResponse(page)

}

// ViewWithList calls the passed view with list and options
func (c *CouchCandy) ViewWithList(ddoc, list, view string, options Options) (*ViewResponse, error) {

	url := fmt.Sprintf("%s/_design/%s/_list/%s/%s/%s%s", createDatabaseURL(c.Session), ddoc, list, ddoc, view, toQueryString(options))
	fmt.Printf("CouchCandy.CallView(%s)\n", url)
	page, err := readJSON(url, c.Get)
	if err != nil {
		return nil, err
	}
	return toViewResponse(page)

}
