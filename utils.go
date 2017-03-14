package couchcandy

import "encoding/json"

// ToString : outputs the passed object as a JSON string.
func ToString(v interface{}) string {
	out, _ := json.Marshal(v)
	return string(out)
}
