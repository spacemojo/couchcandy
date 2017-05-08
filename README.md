**couchcandy**

*Go client for Apache CouchDB* 

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/spacemojo/couchcandy)](https://goreportcard.com/report/github.com/spacemojo/couchcandy)
[![Build Status](https://travis-ci.org/spacemojo/couchcandy.svg?branch=master)](https://travis-ci.org/spacemojo/couchcandy)
[![codecov](https://codecov.io/gh/spacemojo/couchcandy/branch/master/graph/badge.svg)](https://codecov.io/gh/spacemojo/couchcandy)
[![GoDoc](https://godoc.org/github.com/spacemojo/couchcandy?status.svg)](https://godoc.org/github.com/spacemojo/couchcandy)

This is my first try at a GoLang project, be gentle.

To get started : 

```
client := couchcandy.NewCouchCandy(Session{
    Host:     "http://[HOST_IP]",
    Port:     [PORT],
    Database: "database",
    Username: "username",
    Password: "p@$$w0rD",
})
```

From there you can simply call the methods available in the couchcandy client. 

```
info, err := client.GetDatabaseInfo()
```

The returned info object is structured as such : 

```
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
```
