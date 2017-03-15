package couchcandy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// TestCreateDatabaseURL : Tests the CreateDatabaseURL function.
func TestCreateDatabaseURL(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "test", "gotest")
	expected := "http://127.0.0.1:5984/udb"
	url := CreateDatabaseURL(session, "udb")
	if url != expected {
		t.Fail()
	}

}

func TestGetDatabaseInfo(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"db_name":"udb","doc_count":20682,"doc_del_count":0,"update_seq":211591,"purge_seq":0,"compact_running":false,"disk_size":1210183793,"data_size":32983628,"instance_start_time":0,"disk_format_version":6,"committed_update_seq":211591}`)),
		}
		return response, nil
	}
	info, err := couchcandy.GetDatabaseInfo("udb")

	fmt.Println(info)
	fmt.Println(err)

}
