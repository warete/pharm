package main

import (
	"github.com/warete/pharm/cmd/drugs_import"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "pharm",
		Usage: "pharm",
		Commands: []*cli.Command{
			{
				Name:   "drugs-import",
				Usage:  "start app in drugs import mode",
				Action: drugs_import.Cmd,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
