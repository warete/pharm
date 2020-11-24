package main

import (
	"github.com/warete/pharm/cmd/drugs_import"
	"github.com/warete/pharm/cmd/pharm"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	pharm.Init()

	cliApp := &cli.App{
		Name:  "pharm",
		Usage: "pharm",
		Commands: []*cli.Command{
			{
				Name:   "drugs-import",
				Usage:  "start app in drugs import mode",
				Action: drugs_import.Cmd,
			},
			{
				Name:   "app",
				Usage:  "start app",
				Action: pharm.Cmd,
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
