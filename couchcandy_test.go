package couchcandy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func TestDatabaseInfo(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(`{"db_name":"lendr","doc_count":102,"doc_del_count":0,"update_seq":103,"purge_seq":0,"compact_running":false,"disk_size":106600,"data_size":32561,"instance_start_time":"1535733632163685","disk_format_version":6,"committed_update_seq":103}`)),
		}
		return response, nil
	}

	dbInfo, err := couchcandy.DatabaseInfo()
	if err != nil {
		t.Errorf("An error occured whilst fetching DatabaseInfo : %v", err)
	}
	if dbInfo.DBName != "lendr" {
		t.Fail()
	}

}

func TestDatabaseInfoFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Expected failure")
	}

	_, err := couchcandy.DatabaseInfo()
	if err == nil {
		t.Fail()
	}

}

func TestDocument(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"_id":"053cc05f2ee97a0c91d276c9e700194b","_rev":"3-b96f323b37f19c4d1affddf3db3da9c5","type":"com.lendrapp.beans.UserProfile","shortProfile":{"id":null,"firstname":"Patrick","lastname":"Fitzgerald","email":"brun@email.com","organizationId":"053cc05f2ee97a0c91d276c9e700268f","password":"ee0c9435d5e2a07ceaa8abc829990dd3bdd15b7d6d3b0eaac100984da0841530"},"accountType":"PERSONAL","contacts":[],
				"_revisions":{"start":3,"ids":["b96f323b37f19c4d1affddf3db3da9c5","bdeff0741cc1425e5f5b3829a7a9af2f","c76ae1eb708d6eb68974600995b98b70"]}}`)),
		}
		return response, nil
	}

	profile := &UserProfile{}
	err := couchcandy.Document("053cc05f2ee97a0c91d276c9e700194b", profile, Options{
		Revs: true,
		Rev:  "3-b96f323b37f19c4d1affddf3db3da9c5",
	})
	if err != nil || profile.ID != "053cc05f2ee97a0c91d276c9e700194b" /*|| len(profile.Revisions.IDS) != 3*/ {
		t.Fail()
	}

}

func TestDocumentFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "cards", Username: "admin", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestGetDocumentFailure()")
	}

	profile := &UserProfile{}
	err := couchcandy.Document("053cc05f2ee97a0c91d276c9e700194b", profile, Options{
		Revs: false,
		Rev:  "",
	})
	if err == nil {
		t.Fail()
	}

}

func TestDeleteSuccess(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.DeleteHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"_id":"053cc05f2ee97a0c91d276c9e700194b","_rev":"3-b96f323b37f19c4d1affddf3db3da9c5","ok":true}`)),
		}
		return response, nil
	}

	response, err := couchcandy.Delete("053cc05f2ee97a0c91d276c9e700194b", "3-b96f323b37f19c4d1affddf3db3da9c5")
	if err != nil {
		t.Fail()
	}
	if !response.OK {
		t.Fail()
	}

}

func TestDeleteFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.DeleteHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("an error occured when deleting the document")
	}

	_, err := couchcandy.Delete("053cc05f2ee97a0c91d276c9e700194b", "3-b96f323b37f19c4d1affddf3db3da9c5")
	if err == nil {
		t.Fail()
	}

}

func TestAllDatabases(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`["_replicator","_users","baseball","baseball20170228","elements","lendr","social"]`)),
		}
		return response, nil
	}

	names, err := couchcandy.AllDatabases()
	if err != nil {
		t.Fail()
	}
	fmt.Printf("Database names : %v\n", names)

}

func TestAllDatabasesFailure(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestAllDatabasesFailure()")
	}

	_, err := couchcandy.AllDatabases()
	if err == nil {
		t.Fail()
	}

}

func TestChangeNotificatios(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"results":[
				{"seq":19215,"id":"actama99","changes":[{"rev":"1-e860e99218e7c618f3510c48987d6ff0"}]},
				{"seq":19217,"id":"adairbi99","changes":[{"rev":"1-6482114abc008f6ffab3979597fee898"}]},
				{"seq":19993,"id":"armoubi99","changes":[{"rev":"1-4153c31bb3ae6d8553dab186df2b56a3"}]},
				{"seq":20511,"id":"bancrfr99","changes":[{"rev":"1-a67c987380ff807d66308f28698ff0a3"}]},
				{"seq":21679,"id":"bevinte99","changes":[{"rev":"1-9b613c97ee1e03850307ba3c8c36a206"}]},
				{"seq":21697,"id":"bicke99","changes":[{"rev":"1-ca030dc5ace662abf26a3935a3715218"}]},
				{"seq":22177,"id":"bolesjo99","changes":[{"rev":"1-2db4189915c80ed77fb572b0bbf6c03d"}]},
				{"seq":22923,"id":"bristda99","changes":[{"rev":"1-68b20a0bbfd4abe63d359f2c52ac0e9c"}]}]}`)),
		}
		return response, nil
	}

	changes, _ := couchcandy.ChangeNotifications(Options{
		Style: MainOnly,
	})
	if len(changes.Results) != 8 {
		t.Fail()
	}

}

func TestChangeNotificatiosFailure(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error in TestChangeNotificatiosFailure")
	}

	_, err := couchcandy.ChangeNotifications(Options{
		Style: MainOnly,
	})
	if err == nil {
		t.Fail()
	}

}

func TestAddWithID(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}

	couchcandy := NewCouchCandy(session)
	couchcandy.PutHandler = func(string, string) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"id":"1029384756", "rev":"1-b2b5fcc9f6ca0efcd401b9bc40f539cc", "ok": true}`)),
		}
		return response, nil
	}

	response, _ := couchcandy.AddWithID("1029384756", &ShortProfile{})
	if !response.OK {
		t.Fail()
	}

}

func TestAddWithIDFailure(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}

	couchcandy := NewCouchCandy(session)
	couchcandy.PutHandler = func(string, string) (*http.Response, error) {
		return nil, fmt.Errorf("Deliberate error thrown in TestAddWithIDFailure")
	}

	_, err := couchcandy.AddWithID("1029384756", &ShortProfile{})
	if err == nil {
		t.Fail()
	}

}

func TestAddWithIDMarshalError(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	value := make(chan int)
	_, err := couchcandy.AddWithID("102938", value)
	if err == nil {
		t.Fail()
	}

}

func TestUpdate(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}

	couchcandy := NewCouchCandy(session)
	couchcandy.PutHandler = func(string, string) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"id":"1029384756", "rev":"1-b2b5fcc9f6ca0efcd401b9bc40f539cc", "ok": true}`)),
		}
		return response, nil
	}

	response, _ := couchcandy.Update(&ShortProfile{})
	if !response.OK {
		t.Fail()
	}

}

func TestUpdateFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.PutHandler = func(string, string) (*http.Response, error) {
		return nil, fmt.Errorf("Deliberate error from TestUpdateFailure test")
	}

	_, err := couchcandy.Update(&ShortProfile{})
	if err == nil {
		t.Fail()
	}

}

func TestUpdateMarshalError(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	value := make(chan int)
	_, err := couchcandy.Update(value)
	if err == nil {
		t.Fail()
	}

}

func TestAdd(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.PostHandler = func(string, string) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"id":"1029384756", "rev":"1-b2b5fcc9f6ca0efcd401b9bc40f539cc", "ok": true}`)),
		}
		return response, nil
	}

	response, _ := couchcandy.Add(&ShortProfile{})
	if !response.OK {
		t.Fail()
	}

}

func TestAddFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.PostHandler = func(string, string) (*http.Response, error) {
		return nil, fmt.Errorf("Deliberate error in TestAddFailure")
	}

	_, err := couchcandy.Add(&ShortProfile{})
	if err == nil {
		t.Fail()
	}

}

func TestAddMarshalError(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	value := make(chan int)
	_, err := couchcandy.Add(value)
	if err == nil {
		t.Fail()
	}

}

func TestAddDatabase(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.PutHandler = func(string, string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"ok": true}`)),
		}
		return response, nil
	}

	res, err := couchcandy.AddDatabase("unittestdb")
	if err != nil || !res.OK {
		t.Fail()
	}

}

func TestAddDatabaseFailure(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.PutHandler = func(string, string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestAddDatabaseFailure()")
	}

	_, err := couchcandy.AddDatabase("unittestdb")
	if err == nil {
		t.Fail()
	}

}

func TestDeleteDatabase(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
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

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.DeleteHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestDeleteDatabaseFailure()")
	}

	_, err := couchcandy.DeleteDatabase("unittestdb")
	if err == nil {
		t.Fail()
	}

}

func TestDocumentsByKeys(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "userapi", Username: "test", Password: "pwd",
	})
	couchcandy.PostHandler = func(string, string) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
			"total_rows": 5,
			"offset": 0,
			"rows": [
				{
				"id": "1e469de5e70bef8ba57bb5a507cb3cdd",
				"key": "1e469de5e70bef8ba57bb5a507cb3cdd",
				"value": {
					"rev": "5-6ea4f8158560ab4148be647ba28b8572"
				},
				"doc": {
					"_id": "1e469de5e70bef8ba57bb5a507cb3cdd",
					"_rev": "5-6ea4f8158560ab4148be647ba28b8572",
					"firstname": "Charles",
					"lastname": "Darwin",
					"email": "chuck@email.com",
					"status": 4
				}
				},
				{
				"id": "d902c7a42d5c4780af9d7dd3910953a0",
				"key": "d902c7a42d5c4780af9d7dd3910953a0",
				"value": {
					"rev": "1-ac7c71c07efd52e196e7470c9a75f3d7"
				},
				"doc": {
					"_id": "d902c7a42d5c4780af9d7dd3910953a0",
					"_rev": "1-ac7c71c07efd52e196e7470c9a75f3d7",
					"firstname": "Jack",
					"lastname": "Donnaghy",
					"email": "jackattack@email.com",
					"status": 6
				}
				}
			]
			}`)),
		}
		return response, nil
	}

	allDocuments, err := couchcandy.DocumentsByKeys([]string{"Penn", "Teller"}, Options{IncludeDocs: true, Limit: 10})
	if err != nil {
		t.Fail()
	}
	if len(allDocuments.Rows) != 2 {
		t.Fail()
	}

}

func TestDocumentsByKeysFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "userapi", Username: "test", Password: "pwd",
	})
	couchcandy.PostHandler = func(string, string) (*http.Response, error) {
		return nil, fmt.Errorf("An error occured when fetching documents by keys")
	}

	_, err := couchcandy.DocumentsByKeys([]string{""}, Options{IncludeDocs: true, Limit: 10})
	if err == nil {
		t.Fail()
	}

}

func TestAllDocuments(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"total_rows":20682,"offset":0,"rows":[
			{"id":"ALB01","key":"ALB01","value":{"rev":"1-66a2e993bd32c834f4c2bb655b520c42"}},
			{"id":"ALT","key":"ALT","value":{"rev":"2-70d3ae1a59ab2f5be945881afbf6243d"}},
			{"id":"ALT01","key":"ALT01","value":{"rev":"1-d694d4e89aed0b02828d324b26d430c4"}},
			{"id":"ANA","key":"ANA","value":{"rev":"57-bf879cbfc36d7e97e7ddc578f18d675e"}},
			{"id":"ANA01","key":"ANA01","value":{"rev":"1-583301ed352ec7aaea11793618a2fdec"}}
			]}`)),
		}
		return response, nil
	}

	allDocuments, err := couchcandy.Documents(Options{
		Limit:       5,
		IncludeDocs: false,
		Descending:  false,
	})
	if err != nil {
		t.Fail()
	}

	if len(allDocuments.Rows) != 5 {
		t.Fail()
	}

}

func TestAllDocumentsFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from the TestAllDocumentsFailure test")
	}

	_, err := couchcandy.Documents(Options{
		Descending:  false,
		Limit:       5,
		IncludeDocs: false,
	})
	if err == nil {
		t.Fail()
	}

}

func TestCallMapFunction(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"total_rows":52,"offset":39,"rows":[
			{"id":"899b677618b95aad18fae36b6c000310","key":"spades","value":{"_id":"899b677618b95aad18fae36b6c000310","_rev":"1-628016c8a8a997b20359cd5eec8cfd1b","id":"1-spades","suit":"spades","numericValue":1,"name":"Ace","color":"black"}},
			{"id":"899b677618b95aad18fae36b6c000863","key":"spades","value":{"_id":"899b677618b95aad18fae36b6c000863","_rev":"1-86de6dd4a207c2aa28a12340fb397481","id":"2-spades","suit":"spades","numericValue":2,"name":"Two","color":"black"}},
			{"id":"899b677618b95aad18fae36b6c000f45","key":"spades","value":{"_id":"899b677618b95aad18fae36b6c000f45","_rev":"1-9631dd8c00dc96311982bd31ad162090","id":"3-spades","suit":"spades","numericValue":3,"name":"Three","color":"black"}}
			]}`)),
		}
		return response, nil
	}

	docs, err := couchcandy.View("cards", "by_suit", Options{
		Limit:       3,
		IncludeDocs: true,
	})

	if err != nil {
		t.Fail()
	}
	if len(docs.Rows) != 3 {
		t.Fail()
	}

}

func TestCallMapFunctionFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (*http.Response, error) {
		return nil, fmt.Errorf("an error occured whilst calling the map function")
	}

	docs, err := couchcandy.View("cards", "by_suit", Options{})

	if err == nil {
		t.Fail()
	}
	if docs != nil {
		t.Fail()
	}

}

type MockFailingHTTPClient struct{}

func (m *MockFailingHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("Deliberate error from MockCandyHTTPClient.Do")
}

type MockRunningHTTPClient struct{}

func (m *MockRunningHTTPClient) Do(request *http.Request) (*http.Response, error) {
	response := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(`{"total_rows":20682,"offset":0,"rows":[
		{"id":"ALB01","key":"ALB01","value":{"rev":"1-66a2e993bd32c834f4c2bb655b520c42"}},
		{"id":"ALT","key":"ALT","value":{"rev":"2-70d3ae1a59ab2f5be945881afbf6243d"}},
		{"id":"ALT01","key":"ALT01","value":{"rev":"1-d694d4e89aed0b02828d324b26d430c4"}},
		{"id":"ANA","key":"ANA","value":{"rev":"57-bf879cbfc36d7e97e7ddc578f18d675e"}},
		{"id":"ANA01","key":"ANA01","value":{"rev":"1-583301ed352ec7aaea11793618a2fdec"}}
		]}`)),
	}
	return response, nil
}

func TestDefaultHandlerSuccess(t *testing.T) {

	response, _ := defaultHandler(http.MethodGet, "http://127.0.0.1:5984/dbase", &MockRunningHTTPClient{})
	if response == nil {
		t.Fail()
	}

}

func TestDefaultHandlerWithBodySuccess(t *testing.T) {

	response, _ := defaultHandlerWithBody(http.MethodPost, "http://127.0.0.1:5984/dbase", "This is the body for the post.", &MockRunningHTTPClient{})
	if response == nil {
		t.Fail()
	}

}

func TestDefaultHandlerDoFail(t *testing.T) {

	_, err := defaultHandler(http.MethodGet, "http://127.0.0.1:5984/dbase", &MockFailingHTTPClient{})
	if err == nil {
		t.Fail()
	}

}

func TestDefaultHandlerDoRequestFail(t *testing.T) {

	_, err := defaultHandler("\n", "http://127.0.0.1:5984/dbase", &MockFailingHTTPClient{})
	if err == nil {
		t.Fail()
	}

}

func TestDefaultHandlerWithBodyDoFail(t *testing.T) {

	_, err := defaultHandlerWithBody(http.MethodPost, "http://127.0.0.1:5984/dbase", "Body", &MockFailingHTTPClient{})
	if err == nil {
		t.Fail()
	}

}

func TestDefaultHandlerWithBodyDoRequestFail(t *testing.T) {

	_, err := defaultHandlerWithBody(http.MethodPost, "http://127.0.0.1:5984/dbase", "Body", &MockFailingHTTPClient{})
	if err == nil {
		t.Fail()
	}

}

func TestCreateUpdateURL(t *testing.T) {

	url := createPutDocumentURL(Session{}, "{badBodyFormat}")
	if url != "" {
		t.Fail()
	}

}

func TestSafeMarshalError(t *testing.T) {

	value := make(chan int)
	_, err := safeMarshall(value)
	if err == nil {
		t.Fail()
	}

}

func TestToQueryString(t *testing.T) {

	queryString := toQueryString(Options{
		IncludeDocs: true,
		Key:         "\"serge\"",
	})

	expected := "?descending=false&include_docs=true&reduce=false&key=%22serge%22"
	if queryString != expected {
		t.Errorf("Got [%s] - expected [%s]", queryString, expected)
	}

}

func TestToParameters(t *testing.T) {

	parameters := toParameters(Options{
		IncludeDocs: true,
		Limit:       0,
		Key:         "serge",
	})

	fmt.Printf("Paramerers : %v\n", parameters)
	if parameters[0] != "descending=false" || parameters[1] != "include_docs=true" || parameters[2] != "reduce=false" || parameters[3] != "key=serge" {
		t.Fail()
	}

}
