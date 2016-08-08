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
		InitCommand(),
		CreateIssueCommand(),
		ListIssuesCommand(),
		CommentIssueCommand(),
		ShowIssueCommand(),
		CloseIssueCommand(),
	}

	app.Action = func(c *cli.Context) {
		err := ListIssuesCommand().Run(c)
		if err != nil {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}
