package couchcandy

import "net/http"

const (
	// MainOnly Used when getting notifications
	MainOnly string = "main_only"
	// AllDocs Used when getting notifications
	AllDocs string = "all_docs"
)

// CandyHTTPClient Interface that describes a client that executes an
// http request and produces an http response, and error if any.
type CandyHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CandyDocument Struct for holding a CouchDB document.
// Not supposed to be used directly but is required to construct
// your custom types since all documents in CouchDB have these
// 2 attributes and the potential of having Error, Reason and
// _revisions.
type CandyDocument struct {
	ID        string   `json:"_id,omitempty"`
	REV       string   `json:"_rev,omitempty"`
	Error     string   `json:"error,omitempty"`
	Reason    string   `json:"reason,omitempty"`
	Revisions Revision `json:"_revisions,omitempty"`
}

// Revision The revision struct when calling the get document api with revs.
type Revision struct {
	Start int      `json:"start"`
	IDS   []string `json:"ids"`
}

// CouchCandy Struct that provides all CouchDB's API has to offer.
type CouchCandy struct {
	LclSession    Session
	GetHandler    func(string) (*http.Response, error)
	PostHandler   func(string, string) (*http.Response, error)
	PutHandler    func(string, string) (*http.Response, error)
	DeleteHandler func(string) (*http.Response, error)
}

// Changes The struct returned by the call to get change notifications.
type Changes struct {
	Results []Result `json:"results"`
	LastSeq int      `json:"last_seq"`
}

// Result The struct representing a change result.
type Result struct {
	Seq     int      `json:"seq"`
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

// Change The change itself, mainly a revision change on an id.
type Change struct {
	Rev string `json:"rev"`
}

// Options Options available when querying the database.
// Revs : includes revisions or not
// Rev : fetch a specific revision
// Descending : sorting order for the keys
// Limit : number of returned results
// IncludeDocs : includes the whole document or not
// Style :
type Options struct {
	Revs        bool
	Rev         string
	Descending  bool
	Limit       int
	IncludeDocs bool
	Style       string
}

// NewCouchCandy Returns a new CouchCandy struct initialised with the provided values.
func NewCouchCandy(session Session) *CouchCandy {
	return &CouchCandy{
		LclSession:    session,
		GetHandler:    defaultGetHandler,
		PostHandler:   defaultPostHandler,
		PutHandler:    defaultPutHandler,
		DeleteHandler: defaultDeleteHandler,
	}
}

// DatabaseInfo Fetches basic information about a database.
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

// OperationResponse Format of an operation response when a get is not emitted.
type OperationResponse struct {
	ID     string `json:"id"`
	REV    string `json:"rev"`
	OK     bool   `json:"ok"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// Session holds the connection data for a couchcandy session.
type Session struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// AllDocuments This struct contains the response to the all documents call.
type AllDocuments struct {
	TotalRows int   `json:"total_rows"`
	Offset    int   `json:"offset"`
	Rows      []Row `json:"rows"`
}

// Row This is a row in the array of rows on the AllDocuments struct.
type Row struct {
	ID    string      `json:"id"`
	Key   string      `json:"key"`
	Value Value       `json:"value"`
	Doc   interface{} `json:"doc"`
}

// Value The value returned in rows whilst calling CouchDB's _all_docs service.
type Value struct {
	REV string `json:"rev"`
}
