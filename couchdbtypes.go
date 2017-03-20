package couchcandy

import (
	"net/http"
	"strings"
)

// CandyDocument : struct for holding a CouchCandy document.
// Not supposed to be used directly but is required to construct
// your custom types since all documents in CouchDB have these
// 2 attributes.
type CandyDocument struct {
	ID     string `json:"_id"`
	REV    string `json:"_rev"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// CouchCandy : Struct that provides all CouchDB's API has to offer.
type CouchCandy struct {
	LclSession    Session
	GetHandler    func(string) (*http.Response, error)
	PutHandler    func(string, string) (*http.Response, error)
	DeleteHandler func(string) (*http.Response, error)
}

// NewCouchCandy : Returns a new CouchCandy struct initialised with the provided values.
func NewCouchCandy(session Session) *CouchCandy {
	return &CouchCandy{
		LclSession:    session,
		GetHandler:    http.Get,
		PutHandler:    defaultPutHandler,
		DeleteHandler: defaultDeleteHandler,
	}
}

func defaultPutHandler(url string, body string) (*http.Response, error) {

	request, requestError := http.NewRequest(http.MethodPut, url, strings.NewReader(body))
	if requestError != nil {
		return nil, requestError
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil

}

func defaultDeleteHandler(url string) (*http.Response, error) {

	request, requestError := http.NewRequest(http.MethodDelete, url, nil)
	if requestError != nil {
		return nil, requestError
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil

}

// DatabaseInfo : Fetches basic information about a database.
type DatabaseInfo struct {
	DBName             string `json:"db_name"`
	DocCount           int    `json:"doc_count"`
	DocDelCount        int    `json:"doc_del_count"`
	UpdateSeq          int    `json:"update_seq"`
	PurgeSeq           int    `json:"purge_seq"`
	CompactRunning     bool   `json:"compact_running"`
	DiskSize           int    `json:"disk_size"`
	DataSize           int    `json:"data_size"`
	InstanceStartTime  string `json:"instance_start_time"`
	DiskFormatVersion  int    `json:"disk_format_version"`
	CommittedUpdateSeq int    `json:"committed_update_seq"`
}

// OperationResponse : Format of an operation response when a get is not emitted.
type OperationResponse struct {
	ID  string `json:"id"`
	REV string `json:"rev"`
	OK  bool   `json:"ok"`
}

// Session : holds the connection data for a couchcandy session.
type Session struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// AllDocuments : This struct contains the response to the all documents call.
type AllDocuments struct {
	TotalRows int   `json:"total_rows"`
	Offset    int   `json:"offset"`
	Rows      []Row `json:"rows"`
}

// Row : This is a row in the array of rows on the AllDocuments struc.
type Row struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value Value  `json:"value"`
}

// Value : The value returned in rows whilst calling CouchDB's _all_docs service.
type Value struct {
	REV string `json:"rev"`
}

// NewSession : creates a new session initialized with the passed values
func NewSession(host string, port int, database string, username string, password string) Session {
	return Session{
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}
}
