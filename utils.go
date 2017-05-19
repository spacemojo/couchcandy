package couchcandy

import (
	"bytes"
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

func toViewResponse(page []byte) (*ViewResponse, error) {
	viewResponse := &ViewResponse{}
	unmarshallError := json.Unmarshal(page, viewResponse)
	return viewResponse, unmarshallError
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

func toQueryString(options Options) string {

	parameters := toParameters(options)
	buffer := bytes.NewBuffer(make([]byte, 0))
	buffer.WriteString("?")

	for _, parameter := range parameters {
		buffer.WriteString(parameter)
		buffer.WriteString("&")
	}

	return buffer.String()[0 : buffer.Len()-1]

}

func toParameters(options Options) []string {

	parameters := make([]string, 0)

	parameters = append(parameters, fmt.Sprintf("descending=%v", options.Descending))
	if !options.Reduce {
		parameters = append(parameters, fmt.Sprintf("include_docs=%v", options.IncludeDocs))
	}
	parameters = append(parameters, fmt.Sprintf("reduce=%v", options.Reduce))
	if options.Limit != 0 {
		parameters = append(parameters, fmt.Sprintf("limit=%v", options.Limit))
	}
	if options.Key != "" {
		parameters = append(parameters, fmt.Sprintf("key=%s", options.Key))
	}
	if options.StartKey != "" {
		parameters = append(parameters, fmt.Sprintf("start_key=%s", options.StartKey))
	}
	if options.EndKey != "" {
		parameters = append(parameters, fmt.Sprintf("end_key=%s", options.EndKey))
	}
	if options.GroupLevel != 0 {
		parameters = append(parameters, fmt.Sprintf("group_level=%v", options.GroupLevel))
	}
	return parameters

}
