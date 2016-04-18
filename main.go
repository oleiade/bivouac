package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "groundcontrol"
	app.Usage = "make an explosive entrance"
	app.Commands = []cli.Command{
		CreateIssueCommand(),
		ListIssuesCommand(),
		CommentIssueCommand(),
		ShowIssueCommand(),
		CloseIssueCommand(),
	}

	app.Action = func(c *cli.Context) {
	}

	app.Run(os.Args)
}
