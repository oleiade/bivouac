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
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "status",
				Value: "",
				Usage: "filter issues shown based on their status open/closed",
			},
		},
		Action: func(c *cli.Context) {
			store, err := GetOrCreateStore(".groundcontrol")
			if err != nil {
				log.Fatal(err)
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)

			fmt.Fprintln(w, "status\tid\ttitle\topened on\tcomments")

			if c.String("status") != "" {
				issueStatus := ISSUE_OPENED
				if c.String("status") == "closed" {
					issueStatus = ISSUE_CLOSED
				}

				for _, issue := range store.FilterIssues(issueStatus) {
					issueStatusString := "︎!"
					if issue.Status == ISSUE_CLOSED {
						issueStatusString = "✓"
					}

					fmt.Fprintf(
						w, "%s\t#%d\t%s\t%s\t%d\n",
						issueStatusString,
						issue.Id,
						issue.Title,
						issue.Comments[0].CreatedAt.Format("Jan 2 2006"),
						len(issue.Comments)-1,
					)
				}
			} else {
				for _, issue := range store.ListIssues() {
					issueStatusString := "︎!"
					if issue.Status == ISSUE_CLOSED {
						issueStatusString = "✓"
					}

					fmt.Fprintf(
						w, "%s\t#%d\t%s\t%s\t%d\n",
						issueStatusString,
						issue.Id,
						issue.Title,
						issue.Comments[0].CreatedAt.Format("Jan 2 2006"),
						len(issue.Comments)-1,
					)
				}
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

			issueStatus := "!"
			if issue.Status == ISSUE_CLOSED {
				issueStatus = "✓"
			}

			fmt.Printf("%s %s\n", issueStatus, issue.Title)
			fmt.Printf("#%d opened on %s\n", issue.Id, issue.Comments[0].CreatedAt.Format("Jan 2 2006 15:04"))
			fmt.Printf("-----\n")
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
					*NewComment(time.Now(), comment),
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
			cli.BoolFlag{
				Name:  "close",
				Usage: "comment and close the issue",
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

			if c.Bool("close") {
				issue.Close()
			}

			store.Write()
		},
	}
}
