package couchcandy

import "testing"

// TestCreateDatabaseURL : Tests the CreateDatabaseURL function.
func TestCreateDatabaseURL(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "test", "gotest")
	expected := "http://127.0.0.1:5984/udb"
	url := CreateDatabaseURL(session, "udb")
	if url != expected {
		t.Fail()
	}

}
