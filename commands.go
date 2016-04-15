package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/codegangsta/cli"
)

func ListIssuesCommand() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list issues",
		Action: func(c *cli.Context) {
			store, err := GetOrCreateStore(".groundcontrol")
			if err != nil {
				log.Fatal(err)
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)

			fmt.Fprintln(w, "id\ttitle\tcoms\tbody")
			for _, issue := range store.ListIssues() {
				fmt.Fprintf(w, "#%d\t%s\t%d\t%s\n", issue.Id, issue.Title, len(issue.Comments)-1, issue.Comments[0].Body)
			}
			w.Flush()
		},
	}
}

func ShowIssueCommand() cli.Command {
	return cli.Command{
		Name:    "show",
		Aliases: []string{"s"},
		Usage:   "show issue",
		Action: func(c *cli.Context) {
			var err error

			if len(c.Args()) == 0 {
				log.Fatal("Please provide the issue to comment id")
			}

			store, err := GetOrCreateStore(".groundcontrol")
			if err != nil {
				log.Fatal(err)
			}

			id, err := strconv.Atoi(c.Args()[0][1:])
			if err != nil {
				log.Fatal(err)
			}

			issue, err := store.GetIssue(uint(id))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("#%d %s\n", issue.Id, issue.Title)
			fmt.Printf("---\n")
			fmt.Printf("%s\n\n", issue.Comments[0].Body)

			for _, comment := range issue.Comments[1:] {
				fmt.Printf("|\n|\n\n")
				fmt.Printf("%s\n%s\n\n", comment.CreatedAt, comment.Body)
			}
		},
	}
}

func CreateIssueCommand() cli.Command {
	return cli.Command{
		Name:    "add",
		Aliases: []string{"c"},
		Usage:   "create a new issue",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "title",
				Usage: "set the issue tile",
			},
			cli.StringFlag{
				Name:  "comment",
				Usage: "leave a comment",
			},
		},
		Action: func(c *cli.Context) {
			var reader *bufio.Reader
			var title string
			var comment string
			var err error

			store, err := GetOrCreateStore(".groundcontrol")
			if err != nil {
				log.Fatal(err)
			}

			if c.String("title") == "" && c.String("comment") == "" {
				reader = bufio.NewReader(os.Stdin)
				fmt.Print("Title: ")
				title, err = reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}

				fmt.Print("Comment: ")
				comment, err = reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
			} else {
				title = c.String("title")
				comment = c.String("comment")
			}

			issue := NewIssue(
				store.getNextId(),
				title,
				[]Comment{
					*NewComment(time.Now().String(), comment),
				},
			)
			store.AddIssue(*issue)
			store.Write()
		},
	}
}

func CommentIssueCommand() cli.Command {
	return cli.Command{
		Name:    "comment",
		Aliases: []string{},
		Usage:   "comment on a issue",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "comment",
				Usage: "leave a comment",
			},
		},
		Action: func(c *cli.Context) {
			var reader *bufio.Reader
			var issue *Issue
			var comment string
			var err error

			if len(c.Args()) == 0 {
				log.Fatal("Please provide the issue to comment id")
			}

			store, err := GetOrCreateStore(".groundcontrol")
			if err != nil {
				log.Fatal(err)
			}

			if c.String("comment") == "" {
				reader = bufio.NewReader(os.Stdin)

				fmt.Print("Comment: ")
				comment, err = reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
			} else {
				comment = c.String("comment")
			}

			id, err := strconv.Atoi(c.Args()[0][1:])
			if err != nil {
				log.Fatal(err)
			}

			issue, err = store.GetIssue(uint(id))
			if err != nil {
				log.Fatal(err)
			}

			issue.Comment(comment)
			store.Write()
		},
	}
}
