package validator

import "strings"

func removeTopStruct(fields map[string]string) map[string]string {
	resp := map[string]string{}
	for field, err := range fields {
		resp[field[strings.Index(field, ".")+1:]] = err
	}
	return resp
}
