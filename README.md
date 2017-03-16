**couchcandy**

*Go client for Apache CouchDB* 

![Go report card for couchcandy](https://goreportcard.com/badge/github.com/spacemojo/couchcandy)

This is my first try at a GoLang project, be gentle.

To get started : 

~~~~
session := couchcandy.NewSession("http://[HOST_IP]", [PORT], "database", "username", "p@$$w0rD")
client := couchcandy.NewCouchCandy(session)
~~~~

From there you can simply call the methods available in the couchcandy client. 

~~~~
info, err := client.GetDatabaseInfo()
~~~~

The returned info object is structured as such : 

~~~~
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
~~~~
