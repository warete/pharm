package main

import "strings"

func GetQueryStringFromMap(data map[string]string) string {
	elements := make([]string, len(data))

	for key, val := range data {
		elements = append(elements, key+"="+val)
	}

	return strings.Join(elements, "&")
}
