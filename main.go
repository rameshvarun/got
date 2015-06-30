package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

// The file to  be used as an object store
const DBName string = ".got.db"

// The name of the bucket containing OBJECTS
var OBJECTS []byte

// The name of the INFO db Bucket
var INFO []byte

// The key that should always point to the current revision
// that the working copy is based on.
var CURRENT []byte

// This key should point to a list of all heads in the repository
var HEADS []byte

// TODO: Add function that finds root of repository

func init() {
	INFO = []byte("INFO")
	CURRENT = []byte("CURRENT")
	HEADS = []byte("HEADS")
	OBJECTS = []byte("OBJECTS")
}

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

// IgnorePath returns true if a file/directory should be ignored
func IgnorePath(path string) bool {
	if strings.HasPrefix(path, ".") {
		return true
	}
	return false
}

var debug = true

func DebugLog(a string) {
	if debug {
		fmt.Println(a)
	}
}

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
				if _, err := os.Stat(DBName); err == nil {
					fmt.Printf("Got repository already exists in this folder.\n")
				} else {
					db, err := bolt.Open(DBName, 0600, nil)
					if err != nil {
						log.Fatal("Could not create " + DBName + "\n")
					}

					db.Update(func(tx *bolt.Tx) error {
						tx.CreateBucket(INFO)
						tx.CreateBucket(OBJECTS)
						return nil
					})

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

				// Perform operations in a read-only lock
				db.View(func(tx *bolt.Tx) error {
					b := tx.Bucket(INFO)
					heads := b.Get(HEADS)
					if heads != nil {
					}
					fmt.Println("There are no commits in this repository...")
					return nil
				})
			},
		},
		{
			Name:   "status",
			Usage:  "Diff the current working directory with the last commit.",
			Action: Status,
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
			Action: Commit,
		},
		{
			Name:  "checkout",
			Usage: "Checkout a the working tree to a commit.",
			Action: func(c *cli.Context) {
				db := openDB()
				defer db.Close()
			},
		},
		{
			Name:  "revert",
			Usage: "Revert a file to the last commit.",
			Action: func(c *cli.Context) {
				db := openDB()
				defer db.Close()
			},
		},
	}

	app.Run(os.Args)
}
