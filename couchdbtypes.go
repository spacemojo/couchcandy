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
	ID          string                `json:"_id,omitempty"`
	REV         string                `json:"_rev,omitempty"`
	Error       string                `json:"error,omitempty"`
	Reason      string                `json:"reason,omitempty"`
	Attachments map[string]Attachment `json:"_attachments"`
	// Revisions Revision `json:"_revisions,omitempty"`
}

// Attachment is an attachment to a document
type Attachment struct {
	ContentType string `json:"content_type"`
	Revpos      int    `json:"revpos"`
	Digest      string `json:"digest"`
	Length      int    `json:"length"`
	Stub        bool   `json:"stub"`
}

// Revision The revision struct when calling the get document api with revs.
type Revision struct {
	Start int      `json:"start,omitempty"`
	IDS   []string `json:"ids,omitempty"`
}

// CouchCandy Struct that provides all CouchDB's API has to offer.
type CouchCandy struct {
	Session       Session
	GetHandler    func(string) (*http.Response, error)
	PostHandler   func(string, string) (*http.Response, error)
	PutHandler    func(string, string) (*http.Response, error)
	DeleteHandler func(string) (*http.Response, error)
}

// Changes The struct returned by the call to get change notifications.
type Changes struct {
	Results []Result `json:"results,omitempty"`
	LastSeq int      `json:"last_seq,omitempty"`
}

// Result The struct representing a change result.
type Result struct {
	Seq     int      `json:"seq,omitempty"`
	ID      string   `json:"id,omitempty"`
	Changes []Change `json:"changes,omitempty"`
}

// Change The change itself, mainly a revision change on an id.
type Change struct {
	Rev string `json:"rev,omitempty"`
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
	Key         string
	StartKey    string
	EndKey      string
	Reduce      bool
	GroupLevel  int
}

// NewCouchCandy Returns a new CouchCandy struct initialised with the provided values.
func NewCouchCandy(session Session) *CouchCandy {
	return &CouchCandy{
		Session:       session,
		GetHandler:    defaultGetHandler,
		PostHandler:   defaultPostHandler,
		PutHandler:    defaultPutHandler,
		DeleteHandler: defaultDeleteHandler,
	}
}

// DatabaseInfo Fetches basic information about a database.
type DatabaseInfo struct {
	DBName             string `json:"db_name,omitempty"`
	DocCount           int    `json:"doc_count,omitempty"`
	DocDelCount        int    `json:"doc_del_count,omitempty"`
	UpdateSeq          int    `json:"update_seq,omitempty"`
	PurgeSeq           int    `json:"purge_seq,omitempty"`
	CompactRunning     bool   `json:"compact_running,omitempty"`
	DiskSize           int    `json:"disk_size,omitempty"`
	DataSize           int    `json:"data_size,omitempty"`
	InstanceStartTime  string `json:"instance_start_time,omitempty"`
	DiskFormatVersion  int    `json:"disk_format_version,omitempty"`
	CommittedUpdateSeq int    `json:"committed_update_seq,omitempty"`
}

// OperationResponse Format of an operation response when a get is not emitted.
type OperationResponse struct {
	ID     string `json:"id,omitempty"`
	REV    string `json:"rev,omitempty"`
	OK     bool   `json:"ok,omitempty"`
	Error  string `json:"error,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// Session holds the connection data for a couchcandy session.
type Session struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// ViewResponse represents the response sent when a view is called
type ViewResponse struct {
	TotalRows int       `json:"total_rows,omitempty"`
	Offset    int       `json:"offset,omitempty"`
	Error     string    `json:"error,omitempty"`
	Reason    string    `json:"reason,omitempty"`
	Rows      []ViewRow `json:"rows,omitempty"`
}

// ViewRow represents a row in the ViewResponse
type ViewRow struct {
	ID    string      `json:"id,omitempty"`
	Key   interface{} `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// AllDocuments This struct contains the response to the all documents call.
type AllDocuments struct {
	TotalRows int    `json:"total_rows,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Rows      []Row  `json:"rows,omitempty"`
	Error     string `json:"error,omitempty"`
	Reason    string `json:"reason,omitmepty"`
}

// Row This is a row in the array of rows on the AllDocuments struct.
type Row struct {
	ID    string      `json:"id,omitempty"`
	Key   string      `json:"key,omitempty"`
	Value Value       `json:"value,omitempty"`
	Doc   interface{} `json:"doc,omitempty"`
}

// Value The value returned in rows whilst calling CouchDB's _all_docs service.
type Value struct {
	REV string `json:"rev,omitempty"`
}

// AllDocumentsKeys Is used when fetching documents by keys. This struct is passed
// as a POST parameter.
type AllDocumentsKeys struct {
	Keys []string `json:"keys,omitempty"`
}
