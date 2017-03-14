package couchcandy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func getDatabaseURL(session Session, db string) string {
	var buffer bytes.Buffer
	buffer.WriteString(session.Host)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(session.Port))
	buffer.WriteString("/")
	buffer.WriteString(db)
	return buffer.String()
}

// GetDatabaseInfo returns basic information about the passed database.
func GetDatabaseInfo(session Session, db string) (*DatabaseInfo, error) {

	url := getDatabaseURL(session, db)
	res, geterr := http.Get(url)
	if geterr != nil {
		return nil, geterr
	}

	page, geterr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if geterr != nil {
		return nil, geterr
	}

	dbInfo := &DatabaseInfo{}
	json.Unmarshal(page, dbInfo)
	return dbInfo, nil

}
