package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
	"github.com/rameshvarun/got/types"
	"github.com/rameshvarun/got/util"
)

// Status implements the `got status` command.
func Status(c *cli.Context) {
	// Open the database
	db := util.OpenDB()
	defer db.Close()

	// Perform operations in a read-only lock
	err := db.View(func(tx *bolt.Tx) error {
		// Get the current commit sha
		info := tx.Bucket(util.INFO)
		objects := tx.Bucket(util.OBJECTS)

		current := info.Get(util.CURRENT)

		differences := []Difference{}
		if current != nil {
			// Load commit object
			commit := types.DeserializeCommitObject(objects.Get(current))
			util.DebugLog("Comparing working directory to commit '" + commit.Message + "'.")
			differences = TreeDiff(objects, commit.Tree, ".")
		} else {
			// Compare directory to the empty hash
			util.DebugLog("Comparing working directory to empty tree.")
			differences = TreeDiff(objects, types.EMPTY, ".")
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

// Difference represents a modification between two directory trees.
type Difference struct {
	Type     string
	FilePath string
}

// TreeDiff lists the differences between a Tree object in a snapshot
// and a filesystem path.
func TreeDiff(objects *bolt.Bucket, treeHash types.Hash, dir string) []Difference {
	differences := []Difference{}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Could not list files in dir %s: %v", dir, err)
	}
	for _, file := range files {
		if util.IgnorePath(file.Name()) {
			continue
		}

		if file.IsDir() {
			if treeHash.Equal(types.EMPTY) {
				differences = append(differences, TreeDiff(objects, types.EMPTY, path.Join(dir, file.Name()))...)
			}
		} else {
			if treeHash.Equal(types.EMPTY) {
				differences = append(differences, Difference{
					Type:     "A",
					FilePath: path.Join(dir, file.Name()),
				})
			} else {
				treeObject := types.DeserializeTreeObject(objects.Get(treeHash))
				if treeObject.HasFile(file.Name()) {
				}
			}
		}
	}
	return differences
}
