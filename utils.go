package couchcandy

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func createBaseURL(session Session) string {
	return fmt.Sprintf("http://%s:%s@%s:%s", session.Username, session.Password, session.Host[7:], strconv.Itoa(session.Port))
}

func createDatabaseURL(session Session) string {
	return fmt.Sprintf("%s/%s", createBaseURL(session), session.Database)
}

func createPutDocumentURL(session Session, body string) string {
	candyDoc, toCandyError := toCandyDocument(body)
	if toCandyError != nil {
		return ""
	}
	return createDocumentURL(session, candyDoc.ID)
}

func createDocumentURL(session Session, id string) string {
	return fmt.Sprintf("%s/%s", createDatabaseURL(session), id)
}

func createDocumentURLWithOptions(session Session, id string, options Options) string {
	if options.Rev != "" {
		return fmt.Sprintf("%s/?revs=%v&rev=%v", createDocumentURL(session, id), options.Revs, options.Rev)
	}
	return fmt.Sprintf("%s/?revs=%v", createDocumentURL(session, id), options.Revs)
}

func createAllDatabasesURL(session Session) string {
	return fmt.Sprintf("%s/_all_dbs", createBaseURL(session))
}

func toCandyDocument(str string) (*CandyDocument, error) {
	doc := &CandyDocument{}
	err := json.Unmarshal([]byte(str), doc)
	return doc, err
}

func checkOptionsForAllDocuments(options *Options) {
	if options.Limit == 0 {
		options.Limit = 10
	}
}

func toAllDocuments(page []byte) (*AllDocuments, error) {
	allDocuments := &AllDocuments{}
	unmarshallError := json.Unmarshal(page, allDocuments)
	return allDocuments, unmarshallError
}

func toOperationResponse(page []byte) (*OperationResponse, error) {
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
