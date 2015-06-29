package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

// Status implements the `got status` command.
func Status(c *cli.Context) {
	// Open the database
	db := openDB()
	defer db.Close()

	// Perform operations in a read-only lock
	err := db.View(func(tx *bolt.Tx) error {
		// Get the current commit sha
		b := tx.Bucket(INFO)
		current := b.Get(CURRENT)

		differences := []Difference{}
		if current != nil {
			log.Panicln("Not implemented yet")
		} else {
			// Compare directory to the empty hash
			differences = TreeDiff(EMPTY, ".")
		}

		// Print out the found differences
		for _, difference := range differences {
			fmt.Printf("%s %s\n", difference.Type, difference.FilePath)
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error reading from the database.")
	}
}

// TreeDiff lists the differences between a Tree object in a snapshot
// and a filesystem path.
func TreeDiff(tree Hash, dir string) []Difference {
	differences := []Difference{}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Could not list files in dir %s: %v", dir, err)
	}
	for _, file := range files {
		if IgnorePath(file.Name()) {
			continue
		}

		if file.IsDir() {
			if tree == EMPTY {
				differences = append(differences, TreeDiff(EMPTY, path.Join(dir, file.Name()))...)
			}
		} else {
			if tree == EMPTY {
				differences = append(differences, Difference{
					Type:     "A",
					FilePath: path.Join(dir, file.Name()),
				})
			}
		}
	}
	return differences
}
