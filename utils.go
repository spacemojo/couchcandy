package couchcandy

import (
	"fmt"
	"strconv"
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
