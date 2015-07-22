package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
	"github.com/rameshvarun/got/types"
	"github.com/rameshvarun/got/util"
)

// Status implements the `got status` command.
func Status(c *cli.Context) {
	// Open the database
	db := util.OpenDB()
	defer db.Close()

	// Colored output
	yellow := ansi.ColorFunc("yellow+h:black")
	green := ansi.ColorFunc("green+h:black")
	red := ansi.ColorFunc("red+h:black")

	// Perform operations in a read-only lock
	err := db.View(func(tx *bolt.Tx) error {
		// Get the current commit sha
		info := tx.Bucket(util.INFO)
		objects := tx.Bucket(util.OBJECTS)
		current := info.Get(util.CURRENT)

		// Find the differences between the working directory and the tree of the current commit.
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
			line := fmt.Sprintf("%s %s", difference.Type, difference.FilePath)
			if difference.Type == "A" {
				fmt.Println(green(line))
			}
			if difference.Type == "R" {
				fmt.Println(red(line))
			}
			if difference.Type == "M" {
				fmt.Println(yellow(line))
			}
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

func listHasFile(files []os.FileInfo, fileName string) bool {
	for _, file := range files {
		if file.Name() == fileName {
			return true
		}
	}
	return false
}

// TreeDiff lists the differences between a Tree object in a snapshot
// and a filesystem path.
func TreeDiff(objects *bolt.Bucket, treeHash types.Hash, dir string) []Difference {
	differences := []Difference{}

	// Try to list all of the files in this directory.
	files, listErr := ioutil.ReadDir(dir)
	// Try to load in the tree object.
	var treeObject *types.TreeObject = nil
	if !treeHash.Equal(types.EMPTY) {
		treeObject = types.DeserializeTreeObject(objects.Get(treeHash))
	}

	// For each file in the current directory, determine if the file was either added or modified.
	if listErr == nil {
		for _, file := range files {
			if util.IgnorePath(file.Name()) {
				continue
			}

			if file.IsDir() {
				if treeObject == nil {
					differences = append(differences, TreeDiff(objects, types.EMPTY, path.Join(dir, file.Name()))...)
				} else {
					if treeObject.HasFile(file.Name()) {
						differences = append(differences, TreeDiff(objects, treeObject.GetFile(file.Name()), path.Join(dir, file.Name()))...)
					} else {
						differences = append(differences, TreeDiff(objects, types.EMPTY, path.Join(dir, file.Name()))...)
					}
				}
			} else {
				if treeObject == nil {
					differences = append(differences, Difference{
						Type:     "A",
						FilePath: path.Join(dir, file.Name()),
					})
				} else {
					if treeObject.HasFile(file.Name()) {
						fileBytes, err := ioutil.ReadFile(path.Join(dir, file.Name()))
						if err != nil {
							panic(err)
						}
						if !types.CalculateHash(fileBytes).Equal(treeObject.GetFile(file.Name())) {
							differences = append(differences, Difference{
								Type:     "M",
								FilePath: path.Join(dir, file.Name()),
							})
						}
					} else {
						differences = append(differences, Difference{
							Type:     "A",
							FilePath: path.Join(dir, file.Name()),
						})
					}
				}
			}
		}
	}

	// For each file in the the tree object, see if that file was removed in the working directory
	if treeObject != nil {
		for _, entry := range treeObject.Files {
			if !listHasFile(files, entry.Name) {
				if entry.IsDir {
					differences = append(differences, TreeDiff(objects, entry.Hash, path.Join(dir, entry.Name))...)

				} else {
					differences = append(differences, Difference{
						Type:     "R",
						FilePath: path.Join(dir, entry.Name),
					})

				}
			}
		}
	}

	return differences
}
