package couchcandy

import (
	"encoding/json"
	"net/http"
	"os"
)

const (
	// MainOnly Used when getting notifications
	MainOnly string = "main_only"
	// AllDocs Used when getting notifications
	AllDocs string = "all_docs"
	// HeaderContentType is the Content-Type header
	HeaderContentType string = "Content-Type"
	// JSONContentType is the "application/json" content type
	JSONContentType string = "application/json"
)

// CandyHTTPClient Interface that describes a client that executes an
// http request and produces an http response, and error if any.
type CandyHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CandyDocument Struct for holding a CouchDB document.
// Not supposed to be used directly but is required to construct
// your custom types since all documents in CouchDB have these
// 2 attributes and the potential of having Error, Reason,
// Attachments and _revisions.
type CandyDocument struct {
	ID          string                `json:"_id,omitempty"`
	REV         string                `json:"_rev,omitempty"`
	Error       string                `json:"error,omitempty"`
	Reason      string                `json:"reason,omitempty"`
	Attachments map[string]Attachment `json:"_attachments,omitempty"`
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

// NewDBSession returns an inited session with the available environment variables
func NewDBSession() Session {
	return Session{
		Host:     os.Getenv("dbhost"),
		Port:     5984,
		Database: os.Getenv("dbname"),
		Username: os.Getenv("dbusername"),
		Password: os.Getenv("dbpassword"),
	}
}

// NewClient returns an initialized client for connecting to the database
func NewClient() *CouchCandy {
	return NewCouchCandy(NewDBSession())
}

// CouchCandy Struct that provides all CouchDB's API has to offer.
type CouchCandy struct {
	Session  Session
	Get      func(string) (*http.Response, error)
	PostJSON func(string, string) (*http.Response, error)
	PutJSON  func(string, string) (*http.Response, error)
	PutBytes func(string, string, []byte) (*http.Response, error)
	Delete   func(string) (*http.Response, error)
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
	Keys        string
	StartKey    string
	EndKey      string
	Reduce      bool
	GroupLevel  int
	Skip        int
}

// NewCouchCandy Returns a new CouchCandy struct initialised with the provided values.
func NewCouchCandy(session Session) *CouchCandy {
	return &CouchCandy{
		Session:  session,
		Get:      defaultGet,
		PostJSON: defaultPostJSON,
		PutJSON:  defaultPutJSON,
		PutBytes: defaultPutBytes,
		Delete:   defaultDelete,
	}
}

// DatabaseInfo Fetches basic information about a database.
type DatabaseInfo struct {
	DBName             string `json:"db_name,omitempty"`
	DocCount           int    `json:"doc_count,omitempty"`
	DocDelCount        int    `json:"doc_del_count,omitempty"`
	UpdateSeq          string `json:"update_seq,omitempty"`
	PurgeSeq           string `json:"purge_seq,omitempty"`
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

// DesignDocs
type DesignDocs struct {
	MapReduce []*MapReduceDesignDoc `json:"mapreduce"`
	Index     []*IndexDesignDoc     `json:"index"`
}

func NewDesignDocs() *DesignDocs {
	d := &DesignDocs{}
	d.MapReduce = make([]*MapReduceDesignDoc, 0)
	d.Index = make([]*IndexDesignDoc, 0)
	return d
}

type partialDesignDoc struct {
	Language string `json:"language,omitempty"`
}

type MapReduceDesignDoc struct {
	ID       string          `json:"_id,omitempty"`
	REV      string          `json:"_rev,omitempty"`
	Language string          `json:"language,omitempty"`
	Views    map[string]View `json:"views,omitempty"`
}

// Map / Reduce View when the language is "javascript"
type View struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

type IndexDesignDoc struct {
	ID       string               `json:"_id,omitempty"`
	REV      string               `json:"_rev,omitempty"`
	Language string               `json:"language,omitempty"`
	Views    map[string]IndexView `json:"views,omitempty"`
}

// IndexView when the language is "query"
type IndexView struct {
	Map     IndexMap     `json:"map"`
	Reduce  string       `json:"reduce"`
	Options IndexOptions `json:"options"`
}

type IndexMap struct {
	Fields                Fields                 `json:"fields"`
	PartialFilterSelector map[string]interface{} `json:"partial_filter_selector"`
}

type Fields struct {
	CreatedOn string `json:"createdon"`
}

type IndexOptions struct {
	Def Def `json:"def"`
}

type Def struct {
	Fields []string `json:"fields"`
}

// ViewResponse represents the response sent when a view is called
type ViewResponse struct {
	TotalRows int       `json:"total_rows,omitempty"`
	Offset    int       `json:"offset,omitempty"`
	Error     string    `json:"error,omitempty"`
	Reason    string    `json:"reason,omitempty"`
	Rows      []ViewRow `json:"rows,omitempty"`
}

// ViewRow represents a row in the ViewResponse. Both the Key and Value fields are json.RawMessage types
// so that they can be unmarshaled with the desired type in a subsequent step. Since this lib does not have
// any indication as to what will be returned from CouchDB, it is preferred to simply delegate the response
// values to the calling code.
type ViewRow struct {
	ID    string          `json:"id,omitempty"`
	Key   json.RawMessage `json:"key,omitempty"`
	Value json.RawMessage `json:"value,omitempty"`
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
	ID    string          `json:"id,omitempty"`
	Key   string          `json:"key,omitempty"`
	Value Value           `json:"value,omitempty"`
	Doc   json.RawMessage `json:"doc,omitempty"`
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
