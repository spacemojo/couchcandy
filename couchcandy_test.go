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

func TestGetDatabaseInfo(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "udb", Username: "test", Password: "gotest",
	})
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

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "udb", Username: "test", Password: "gotest",
	})
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
	err := couchcandy.GetDocument("053cc05f2ee97a0c91d276c9e700194b", profile, Options{
		Revs: true,
		Rev:  "3-b96f323b37f19c4d1affddf3db3da9c5",
	})
	if err != nil || profile.ID != "053cc05f2ee97a0c91d276c9e700194b" || len(profile.Revisions.IDS) != 3 {
		t.Fail()
	}

}

func TestGetDocumentFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestGetDocumentFailure()")
	}

	profile := &UserProfile{}
	err := couchcandy.GetDocument("053cc05f2ee97a0c91d276c9e700194b", profile, Options{
		Revs: false,
		Rev:  "",
	})
	if err == nil {
		t.Fail()
	}

}

func TestDeleteDocumentSuccess(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.DeleteHandler = func(string) (resp *http.Response, e error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"_id":"053cc05f2ee97a0c91d276c9e700194b","_rev":"3-b96f323b37f19c4d1affddf3db3da9c5","ok":true}`)),
		}
		return response, nil
	}

	response, err := couchcandy.DeleteDocument("053cc05f2ee97a0c91d276c9e700194b", "3-b96f323b37f19c4d1affddf3db3da9c5")
	if err != nil {
		t.Fail()
	}
	if !response.OK {
		t.Fail()
	}

}

func TestGetAllDatabases(t *testing.T) {

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

	names, err := couchcandy.GetAllDatabases()
	if err != nil {
		t.Fail()
	}
	fmt.Printf("Database names : %v\n", names)

}

func TestGetAllDatabasesFailure(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from TestGetAllDatabasesFailure()")
	}

	_, err := couchcandy.GetAllDatabases()
	if err == nil {
		t.Fail()
	}

}

func TestGetChangeNotifications(t *testing.T) {

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

	changes, _ := couchcandy.GetChangeNotifications(Options{
		Style: MainOnly,
	})
	if len(changes.Results) != 8 {
		t.Fail()
	}

}

func TestGetChangeNotificationsFailure(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
	couchcandy := NewCouchCandy(session)
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error in TestGetChangeNotificationsFailure")
	}

	_, err := couchcandy.GetChangeNotifications(Options{
		Style: MainOnly,
	})
	if err == nil {
		t.Fail()
	}

}

func TestPutDocumentWithID(t *testing.T) {

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

	response, _ := couchcandy.PutDocumentWithID("1029384756", &ShortProfile{})
	if !response.OK {
		t.Fail()
	}

}

func TestPutDocumentWithIDFailure(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}

	couchcandy := NewCouchCandy(session)
	couchcandy.PutHandler = func(string, string) (*http.Response, error) {
		return nil, fmt.Errorf("Deliberate error thrown in TestPutDocumentWithIDFailure")
	}

	_, err := couchcandy.PutDocumentWithID("1029384756", &ShortProfile{})
	if err == nil {
		t.Fail()
	}

}

func TestPutDocument(t *testing.T) {

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

	response, _ := couchcandy.PutDocument(&ShortProfile{})
	if !response.OK {
		t.Fail()
	}

}

func TestPutDocumentFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.PutHandler = func(string, string) (*http.Response, error) {
		return nil, fmt.Errorf("Deliberate error from TestPutDocumentFailure test")
	}

	_, err := couchcandy.PutDocument(&ShortProfile{})
	if err == nil {
		t.Fail()
	}

}

func TestPostDocument(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.PostHandler = func(string, string) (*http.Response, error) {
		response := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"id":"1029384756", "rev":"1-b2b5fcc9f6ca0efcd401b9bc40f539cc", "ok": true}`)),
		}
		return response, nil
	}

	response, _ := couchcandy.PostDocument(&ShortProfile{})
	if !response.OK {
		t.Fail()
	}

}

func TestPostDocumentFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.PostHandler = func(string, string) (*http.Response, error) {
		return nil, fmt.Errorf("Deliberate error in TestPostDocumentFailure")
	}

	_, err := couchcandy.PostDocument(&ShortProfile{})
	if err == nil {
		t.Fail()
	}

}

func TestPutDatabase(t *testing.T) {

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

	res, err := couchcandy.PutDatabase("unittestdb")
	if err != nil || !res.OK {
		t.Fail()
	}

}

func TestPutDatabaseFailure(t *testing.T) {

	session := Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	}
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

func TestGetDocumentsByKeys(t *testing.T) {

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

	allDocuments, err := couchcandy.GetDocumentsByKeys([]string{"Penn", "Teller"}, Options{IncludeDocs: true, Limit: 10})
	if err != nil {
		t.Fail()
	}
	if len(allDocuments.Rows) != 2 {
		t.Fail()
	}

}

func TestGetAllDocuments(t *testing.T) {

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

	allDocuments, err := couchcandy.GetAllDocuments(Options{
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

func TestGetAllDocumentsFailure(t *testing.T) {

	couchcandy := NewCouchCandy(Session{
		Host: "http://127.0.0.1", Port: 5984, Database: "lendr", Username: "test", Password: "gotest",
	})
	couchcandy.GetHandler = func(string) (resp *http.Response, e error) {
		return nil, fmt.Errorf("Deliberate error from the TestGetAllDocumentsFailure test")
	}

	_, err := couchcandy.GetAllDocuments(Options{
		Descending:  false,
		Limit:       5,
		IncludeDocs: false,
	})
	if err == nil {
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

func TestCreatePutDocumentURL(t *testing.T) {

	url := createPutDocumentURL(Session{}, "{badBodyFormat}")
	if url != "" {
		t.Fail()
	}

}

func TestCheckOptionsForAllDocuments(t *testing.T) {

	options := &Options{
		Limit: 0,
	}
	checkOptionsForAllDocuments(options)
	if options.Limit != 10 {
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
