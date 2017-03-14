package couchcandy

// CandyDocument : struct for holding a CouchCandy document
type CandyDocument struct {
	ID  string `json:"_id"`
	REV string `json:"_rev"`
}

// NewCandyDocument : creates a new CandyDocument struct
func NewCandyDocument() *CandyDocument {
	document := &CandyDocument{
		ID:  "",
		REV: "",
	}
	return document
}
