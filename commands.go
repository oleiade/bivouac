package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/codegangsta/cli"
	"github.com/kataras/iris"
)

func initCommand() cli.Command {
	return cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "creates an empty bivouac issue tracker",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "project-name",
				Value: "",
				Usage: "set the issue tracker's project name",
			},
		},
		Action: func(c *cli.Context) {
			var reader *bufio.Reader
			var err error

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("fatal: unable to compute current working directory; reason: %v", err)
			}

			storePath := path.Join(cwd, bivouacFile)

			s, err := GetOrCreateStore(storePath)
			if err != nil {
				log.Fatalf("fatal: unable to get or create Bivouac file; reason: %v", err)
			}

			if c.String("project-name") == "" {
				var projectName string
				reader = bufio.NewReader(os.Stdin)
				fmt.Print("Project name: ")
				projectName, err = reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
				s.ProjectName = strings.Trim(projectName, "\n")
			} else {
				s.ProjectName = c.String("project-name")
			}

			err = s.Write()
			if err != nil {
				log.Fatalf("unable to sync store to disk; reason %v", err)
			}

			fmt.Printf("Initialized empty Bivouac issue tracker: %s\n", storePath)
		},
	}
}

func listIssuesCommand() cli.Command {
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
			storePath, err := findBivouacFile()
			if err != nil {
				log.Fatalf("fatal: No bivouac issue tracker found (nor in any of the parent directories): .bivouac")
			}

			store, err := GetOrCreateStore(storePath)
			if err != nil {
				log.Fatal(err)
			}

			ok := store.HasIssues()
			if !ok {
				log.Fatal("fatal: your issue tracker does not have any issues yet")
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)

			fmt.Fprintln(w, "status\tid\ttitle\topened on\tcomments")

			if c.String("status") != "" {
				issueStatus := IssueOpen
				if c.String("status") == "closed" {
					issueStatus = IssueClosed
				}

				for _, issue := range store.FilterIssues(issueStatus) {
					issueStatusString := "︎!"
					if issue.Status == IssueClosed {
						issueStatusString = "✓"
					}

					fmt.Fprintf(
						w, "%s\t#%d\t%s\t%s\t%d\n",
						issueStatusString,
						issue.ID,
						issue.Title,
						time.Unix(issue.CreatedAt, 0).Format("Jan 2 2006"),
						len(issue.Comments),
					)
				}
			} else {
				for _, issue := range store.ListIssues() {
					issueStatusString := "︎!"
					if issue.Status == IssueClosed {
						issueStatusString = "✓"
					}

					fmt.Fprintf(
						w, "%s\t#%d\t%s\t%s\t%d\n",
						issueStatusString,
						issue.ID,
						issue.Title,
						time.Unix(issue.CreatedAt, 0).Format("Jan 2 2006"),
						len(issue.Comments),
					)
				}
			}
			w.Flush()
		},
	}
}

func showIssueCommand() cli.Command {
	return cli.Command{
		Name:    "show",
		Aliases: []string{"s"},
		Usage:   "show issue",
		Action: func(c *cli.Context) {
			var err error

			if len(c.Args()) == 0 {
				log.Fatal("Please provide the issue to comment id")
			}

			storePath, err := findBivouacFile()
			if err != nil {
				log.Fatal(err)
			}

			store, err := GetOrCreateStore(storePath)
			if err != nil {
				log.Fatal(err)
			}

			id, err := strconv.Atoi(c.Args()[0])
			if err != nil {
				log.Fatal(err)
			}

			issue, err := store.GetIssue(uint(id))
			if err != nil {
				log.Fatal(err)
			}

			issueStatus := "!"
			if issue.Status == IssueClosed {
				issueStatus = "✓"
			}

			fmt.Printf("%s %s\n", issueStatus, issue.Title)
			fmt.Printf("#%d opened on %s\n", issue.ID, time.Unix(issue.CreatedAt, 0).Format("Jan 2 2006 15:04"))
			fmt.Printf("-----\n")
			fmt.Printf("%s\n\n", issue.Description)

			for _, comment := range issue.Comments {
				fmt.Printf("\n\n")
				fmt.Printf("Commented on %s\n-----\n%s\n\n", time.Unix(comment.CreatedAt, 0).Format("Jan 2 2006 15:04"), comment.Body)
			}
		},
	}
}

func createIssueCommand() cli.Command {
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
				Name:  "description",
				Usage: "describe the issue",
			},
		},
		Action: func(c *cli.Context) {
			var reader *bufio.Reader
			var title string
			var description string
			var err error

			storePath, err := findBivouacFile()
			if err != nil {
				log.Fatal(err)
			}

			store, err := GetOrCreateStore(storePath)
			if err != nil {
				log.Fatal(err)
			}

			if c.String("title") == "" && c.String("description") == "" {
				reader = bufio.NewReader(os.Stdin)
				fmt.Print("Title: ")
				title, err = reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}

				fmt.Print("Description: ")
				description, err = reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
			} else {
				title = c.String("title")
				description = c.String("description")
			}

			issue := NewIssue(
				store.getNextID(),
				title,
				description,
			)
			store.AddIssue(*issue)
			store.Write()

			fmt.Println(issue.ID)
		},
	}
}

func commentIssueCommand() cli.Command {
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

			storePath, err := findBivouacFile()
			if err != nil {
				log.Fatal(err)
			}

			store, err := GetOrCreateStore(storePath)
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

			id, err := strconv.Atoi(c.Args()[0])
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

func closeIssueCommand() cli.Command {
	return cli.Command{
		Name:    "close",
		Aliases: []string{},
		Usage:   "close an issue",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "comment",
				Value: "",
				Usage: "close with a comment",
			},
		},
		Action: func(c *cli.Context) {
			var issue *Issue
			var err error

			if len(c.Args()) == 0 {
				log.Fatal("Please provide the issue to close")
			}

			storePath, err := findBivouacFile()
			if err != nil {
				log.Fatal(err)
			}

			store, err := GetOrCreateStore(storePath)
			if err != nil {
				log.Fatal(err)
			}

			id, err := strconv.Atoi(c.Args()[0])
			if err != nil {
				log.Fatal(err)
			}

			issue, err = store.GetIssue(uint(id))
			if err != nil {
				log.Fatal(err)
			}

			if c.String("comment") != "" {
				issue.Comment(c.String("comment"))
			}

			issue.Close()
			store.Write()
		},
	}
}

func serveCommand() cli.Command {
	return cli.Command{
		Name:    "serve",
		Aliases: []string{},
		Usage:   "start an issue server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "transport",
				Value: ":8080",
				Usage: "set the server transport",
			},
		},
		Action: func(c *cli.Context) {

			iris.Get("/", getDocument)

			iris.Get("/issues", getIssues)
			iris.Get("/issues/:issue_id", getIssue)

			iris.Get("/issues/:issue_id/comments", getComments)

			// serve requests at http://localhost:8080
			iris.Listen(":8080")
		},
	}
}
