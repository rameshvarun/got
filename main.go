package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/rameshvarun/got/commands"
	"github.com/rameshvarun/got/util"
	"github.com/urfave/cli"
)

// Application entry point
func main() {
	app := cli.NewApp()
	app.Name = "got"
	app.Usage = "A VCS written in golang."

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Create an empty got repository in the current directory.",
			Action: func(c *cli.Context) {
				if _, err := os.Stat(util.DBName); err == nil {
					fmt.Printf("Got repository already exists in this folder.\n")
				} else {
					db, err := bolt.Open(util.DBName, 0600, nil)
					if err != nil {
						log.Fatal("Could not create " + util.DBName + "\n")
					}

					db.Update(func(tx *bolt.Tx) error {
						tx.CreateBucket(util.INFO)
						tx.CreateBucket(util.OBJECTS)
						return nil
					})

					defer db.Close()
				}
			},
		},
		{
			Name:   "log",
			Usage:  "Show commit logs.",
			Action: commands.Log,
		},
		{
			Name:   "status",
			Usage:  "Diff the current working directory with the last commit.",
			Action: commands.Status,
		},
		{
			Name:  "commit",
			Usage: "Record changes to the repository.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "message, m",
					Usage: "Commit message",
				},
				cli.StringFlag{
					Name:  "author, a",
					Usage: "Commit author",
				},
			},
			Action: commands.Commit,
		},
		{
			Name:  "checkout",
			Usage: "Checkout a the working tree to a commit.",
			Action: func(c *cli.Context) {
				db := util.OpenDB()
				defer db.Close()
			},
		},
		{
			Name:   "revert",
			Usage:  "Revert a file to the last commit.",
			Action: commands.Revert,
		},
	}

	app.Run(os.Args)
}
