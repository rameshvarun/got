package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
	"github.com/rameshvarun/got/types"
	"github.com/rameshvarun/got/util"
)

// TreeLookup gets the contents of a certain file path in the given tree
func TreeLookup(objects *bolt.Bucket, tree *types.TreeObject, filepath string) []byte {
	// Split the path into sperate components
	components := strings.Split(filepath, string(os.PathSeparator))
	if !tree.HasFile(components[0]) {
		log.Fatalln("tree object does not have file " + components[0])
	}

	hash := tree.GetFile(components[0])
	if len(components) > 1 {
		// If there is more than one component left, expect the object to be another directory.
		nextTree := types.DeserializeTreeObject(objects.Get(hash))
		return TreeLookup(objects, nextTree, strings.Join(components[1:], string(os.PathSeparator)))
	}

	// Return the file data stored by that hash.
	return objects.Get(hash)
}

// Revert implements got revert
func Revert(c *cli.Context) {
	// Open the database
	db := util.OpenDB()
	defer db.Close()

	// Perform operations in a read-only lock
	err := db.View(func(tx *bolt.Tx) error {
		// Get the current commit sha
		info := tx.Bucket(util.INFO)
		objects := tx.Bucket(util.OBJECTS)
		current := info.Get(util.CURRENT)

		// Get the commit object and associated tree
		commit := types.DeserializeCommitObject(objects.Get(current))
		tree := types.DeserializeTreeObject(objects.Get(commit.Tree))

		for _, file := range c.Args() {
			fmt.Println("Reverting " + file + "...")
			fileData := TreeLookup(objects, tree, path.Clean(file))
			ioutil.WriteFile(file, fileData, 0644)
		}

		return nil
	})

	if err != nil {
		log.Fatal("Error reading from the database.")
	}
}
