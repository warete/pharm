package drugs_import

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/warete/pharm/config"
	"log"
)

func Cmd(c *cli.Context) error {
	AppConfig, err := config.Init(c.String("config"))
	if err != nil {
		log.Fatal(err)
	}

	meImport, err := Init("https://drugs.medelement.com/search/load_data", map[string]string{
		"searched-data":        "drugs",
		"parent_category_code": "",
		"category_code":        "",
		"q":                    "",
		"result-type":          "json",
		"skip":                 "0",
	}, AppConfig)
	if err != nil {
		return err
	}

	fmt.Println(meImport.GetAll())

	//meImport.Run()

	return nil
}
