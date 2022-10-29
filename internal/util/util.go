package util

import "encoding/json"

func PrettyString(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
