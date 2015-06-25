package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

// The file to  be used as an object store
const DBName string = ".got.db"

// Tries to open the database, printing approproate errors if it doesn't exist,
// or could not be opened.
func openDB() *bolt.DB {
	if _, err := os.Stat(DBName); err == nil {
		db, err := bolt.Open(DBName, 0600, nil)
		if err != nil {
			log.Fatal("Could not open " + DBName + "\n")
		}
		return db
	}
	log.Fatal(DBName + " does not exist.\n")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "got"
	app.Usage = "A VCS written in golang."
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Create an empty got repository in the current directory.",
			Action: func(c *cli.Context) {
				if _, err := os.Stat(DBName); err == nil {
					fmt.Printf("Got repository already exists in this folder.\n")
				} else {
					db, err := bolt.Open(DBName, 0600, nil)
					if err != nil {
						log.Fatal("Could not create " + DBName + "\n")
					}
					defer db.Close()
				}
			},
		},
		{
			Name:  "log",
			Usage: "Show commit logs.",
			Action: func(c *cli.Context) {
				db := openDB()
				defer db.Close()
			},
		},
		{
			Name:  "status",
			Usage: "Diff the current working directory with the last commit.",
			Action: func(c *cli.Context) {
				db := openDB()
				defer db.Close()
			},
		},
		{
			Name:  "commit",
			Usage: "Record changes to the repository.",
			Action: func(c *cli.Context) {
				db, err := bolt.Open(DBName, 0600, nil)
				if err != nil {
					log.Fatal("Could not initialize got.db.\n")
				}
				defer db.Close()
			},
		},
		{
			Name:  "checkout",
			Usage: "Checkout a the working tree to a commit.",
			Action: func(c *cli.Context) {
			},
		},
		{
			Name:  "revert",
			Usage: "Revert a file to the last commit.",
			Action: func(c *cli.Context) {
				db, err := bolt.Open(DBName, 0600, nil)
				if err != nil {
					log.Fatal("Could not initialize got.db.\n")
				}
				defer db.Close()
			},
		},
	}

	app.Run(os.Args)
}
