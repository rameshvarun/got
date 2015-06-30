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
		info := tx.Bucket(INFO)
		objects := tx.Bucket(OBJECTS)

		current := info.Get(CURRENT)

		differences := []Difference{}
		if current != nil {
			// Load commit object
			commit := DeserializeCommitObject(objects.Get(current))
			DebugLog("Comparing working directory to commit '" + commit.Message + "'.")
			differences = TreeDiff(objects, commit.Tree, ".")
		} else {
			// Compare directory to the empty hash
			DebugLog("Comparing working directory to empty tree.")
			differences = TreeDiff(objects, EMPTY, ".")
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
func TreeDiff(objects *bolt.Bucket, treeHash Hash, dir string) []Difference {
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
			if treeHash.Equal(EMPTY) {
				differences = append(differences, TreeDiff(objects, EMPTY, path.Join(dir, file.Name()))...)
			}
		} else {
			if treeHash.Equal(EMPTY) {
				differences = append(differences, Difference{
					Type:     "A",
					FilePath: path.Join(dir, file.Name()),
				})
			} else {
				treeObject := DeserializeTreeObject(objects.Get(treeHash))
        if treeObject.HasFile(name string)
			}
		}
	}
	return differences
}
