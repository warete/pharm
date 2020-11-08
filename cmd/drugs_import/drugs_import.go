package drugs_import

import (
	"github.com/urfave/cli/v2"
)

func Cmd(c *cli.Context) error {
	meImport, err := NewImporter("https://drugs.medelement.com/search/load_data", map[string]string{
		"searched-data":        "drugs",
		"parent_category_code": "",
		"category_code":        "",
		"q":                    "",
		"result-type":          "json",
		"skip":                 "0",
	})
	if err != nil {
		return err
	}

	meImport.Run()

	return nil
}
