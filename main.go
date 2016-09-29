package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "bivouac"
	app.Usage = "make an explosive entrance"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		initCommand(),
		createIssueCommand(),
		listIssuesCommand(),
		commentIssueCommand(),
		showIssueCommand(),
		closeIssueCommand(),
		serveCommand(),
	}

	app.Action = func(c *cli.Context) {
		err := listIssuesCommand().Run(c)
		if err != nil {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}
