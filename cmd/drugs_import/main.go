package main

import "log"

var serviceUrl string = "https://drugs.medelement.com/search/load_data"
var serviceQueryParams = map[string]string{
	"searched-data":        "drugs",
	"parent_category_code": "",
	"category_code":        "",
	"q":                    "",
	"result-type":          "json",
	"skip":                 "0",
}

func main() {
	meImport, err := NewImporter("https://drugs.medelement.com/search/load_data", map[string]string{
		"searched-data":        "drugs",
		"parent_category_code": "",
		"category_code":        "",
		"q":                    "",
		"result-type":          "json",
		"skip":                 "0",
	})
	if err != nil {
		log.Fatal(err)
	}

	meImport.Run()
}
