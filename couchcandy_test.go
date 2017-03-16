package couchcandy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type UserProfile struct {
	CandyDocument
	Type        string       `json:"type"`
	AccountType string       `json:"accountType"`
	Short       ShortProfile `json:"shortProfile"`
}

type ShortProfile struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

// TestCreateDatabaseURL : Tests the CreateDatabaseURL function.
func TestCreateDatabaseURL(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "udb", "test", "gotest")
	expected := "http://test:gotest@127.0.0.1:5984/udb"
	url := CreateDatabaseURL(session)
	if url != expected {
		t.Fail()
	}

}

func TestGetDatabaseInfo(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "udb", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"db_name":"udb","doc_count":20682,"doc_del_count":0,"update_seq":211591,"purge_seq":0,"compact_running":false,"disk_size":1210183793,"data_size":32983628,"instance_start_time":"0","disk_format_version":6,"committed_update_seq":211591}`)),
		}
		return response, nil
	}
	info, err := couchcandy.GetDatabaseInfo()

	if info.DocCount != 20682 || info.DBName != "udb" || err != nil {
		t.Fail()
	}

}

func TestGetDatabaseInfoFailure(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "udb", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("%s", "This is a deliberate error in unit tests (TestGetDatabaseInfoFailure)")
	}
	_, err := couchcandy.GetDatabaseInfo()
	if err == nil {
		t.Fail()
	} else {
		fmt.Println(err)
	}

}

func TestGetDocument(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "lendr", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"_id":"053cc05f2ee97a0c91d276c9e700194b","_rev":"3-b96f323b37f19c4d1affddf3db3da9c5","type":"com.lendrapp.beans.UserProfile","shortProfile":{"id":null,"firstname":"Patrick","lastname":"Fitzgerald","email":"brun@email.com","organizationId":"053cc05f2ee97a0c91d276c9e700268f","password":"ee0c9435d5e2a07ceaa8abc829990dd3bdd15b7d6d3b0eaac100984da0841530"},"accountType":"PERSONAL","contacts":[]}`)),
		}
		return response, nil
	}

	profile := &UserProfile{}
	err := couchcandy.GetDocument("053cc05f2ee97a0c91d276c9e700194b", profile)
	if err != nil || profile.ID != "053cc05f2ee97a0c91d276c9e700194b" {
		t.Fail()
	}

}

func TestGetDocumentFailure(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "lendr", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestGetDocumentFailure()")
	}

	profile := &UserProfile{}
	err := couchcandy.GetDocument("053cc05f2ee97a0c91d276c9e700194b", profile)
	if err == nil {
		t.Fail()
	}

}

func TestGetAllDatabases(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "lendr", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`["_replicator","_users","baseball","baseball20170228","elements","lendr","social"]`)),
		}
		return response, nil
	}

	names, err := couchcandy.GetAllDatabases()
	if err != nil {
		t.Fail()
	}
	fmt.Printf("Database names : %v\n", names)

}

func TestGetAllDatabasesFailure(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "lendr", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestGetAllDatabasesFailure()")
	}

	_, err := couchcandy.GetAllDatabases()
	if err == nil {
		t.Fail()
	}

}

func TestPutDatabase(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "lendr", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.PutHandler = func(string, string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"ok": true}`)),
		}
		return response, nil
	}

	res, err := couchcandy.PutDatabase("unittestdb")
	if err != nil || !res.OK {
		t.Fail()
	}

}

func TestPutDatabaseFailure(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "lendr", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.PutHandler = func(string, string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestPutDatabaseFailure()")
	}

	_, err := couchcandy.PutDatabase("unittestdb")
	if err == nil {
		t.Fail()
	}

}

func TestDeleteDatabase(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "lendr", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.DeleteHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"ok": true}`)),
		}
		return response, nil
	}

	res, err := couchcandy.DeleteDatabase("unittestdb")
	if err != nil || !res.OK {
		t.Fail()
	}

}

func TestDeleteDatabaseFailure(t *testing.T) {

	session := NewSession("http://127.0.0.1", 5984, "lendr", "test", "gotest")
	couchcandy := NewCouchCandy(session)
	couchcandy.DeleteHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"ok": true}`)),
		}
		return response, nil
	}

	res, err := couchcandy.DeleteDatabase("unittestdb")
	if err != nil || !res.OK {
		t.Fail()
	}

}
