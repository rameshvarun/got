package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
	"github.com/rameshvarun/got/types"
	"github.com/rameshvarun/got/util"
)

// TakeSnapshot recursivelly snapshots the current directory.
func TakeSnapshot(b *bolt.Bucket, filePath string) types.Hash {
	// Determine if path is a directory or a file
	info, err := os.Stat(filePath)
	if err != nil {
		log.Fatalln("Could not call os.Stat on " + filePath)
	}

	if info.IsDir() {
		// If the path is a directory, we need to create a tree object
		treeObject := types.NewTreeObject()

		// Enumerate all of the files in this directory
		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			log.Fatalf("Could not list files in dir %s: %v", filePath, err)
		}
		for _, file := range files {
			if util.IgnorePath(file.Name()) {
				continue
			}

			// Snapshot the file, and add it to the current tree object
			hash := TakeSnapshot(b, path.Join(filePath, file.Name()))
			treeObject.AddFile(file.Name(), hash)
		}

		// Calculate hash of the directory
		serialized := treeObject.Serialize()
		hash := types.CalculateHash(serialized)

		// Put directory into database, addressed by it's hash
		b.Put(hash, serialized)
		util.DebugLog("Put directory \"" + filePath + "\" into database with hash " + hash.String())
		return hash
	}

	// Calculate hash of file contents
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	hash := types.CalculateHash(file)

	// Put file into database, addressed by it's hash
	b.Put(hash, file)
	util.DebugLog("Put file \"" + filePath + "\" into database with hash " + hash.String())
	return hash
}

// Commit implements `got commit`
func Commit(c *cli.Context) {
	db := util.OpenDB()
	defer db.Close()

	if len(c.String("message")) == 0 {
		log.Fatalln("Must supply a commit message.")
	}

	if len(c.String("author")) == 0 {
		log.Fatalln("Must supply a commit author.")
	}

	// Perform operations in a write lock
	err := db.Update(func(tx *bolt.Tx) error {
		info := tx.Bucket(util.INFO)       // The bucket holding repositry metadata
		objects := tx.Bucket(util.OBJECTS) // The bucket holding got objects
		current := info.Get(util.CURRENT)

		// Create a commit object from the current directory snapshot
		util.DebugLog("Taking snapshot of repo...")
		hash := TakeSnapshot(tx.Bucket(util.OBJECTS), ".")
		util.DebugLog("Repo snapshot has hash " + hash.String() + ".")

		commit := types.CommitObject{
			Author:  c.String("author"),
			Message: c.String("message"),
			Time:    time.Now(),
			Tree:    hash,
			Parents: make([]types.Hash, 0),
		}

		// There is a 'current' commit, then it is the parent of the new commit
		if current != nil {
			commit.Parents = append(commit.Parents, current)
		}

		commitBytes := commit.Serialize()
		commitSha := types.CalculateHash(commitBytes)
		objects.Put(commitSha, commitBytes)
		util.DebugLog("Created commit of sha " + commitSha.String() + ".")

		// CURRENT now corresponds to the commit that we just created
		info.Put(util.CURRENT, commitSha)

		// Update heads
		headsBytes := info.Get(util.HEADS)
		if headsBytes != nil {
			heads := types.DeserializeHashes(headsBytes)
			util.DebugLog("Currently " + strconv.Itoa(len(heads)) + " heads in this repo...")

			// Remove current, if it is a head
			for i, head := range heads {
				if head.Equal(current) {
					heads = append(heads[:i], heads[i+1:]...)
					break
				}
			}

			// Add the new commit as a head
			heads = append(heads, commitSha)

			info.Put(util.HEADS, types.SerializeHashes(heads))
		} else {
			util.DebugLog("Currently no heads in this repo. Putting HEADS: " + fmt.Sprint(([]types.Hash{commitSha})))
			info.Put(util.HEADS, types.SerializeHashes([]types.Hash{commitSha}))
		}

		return nil
	})
	if err != nil {
		log.Fatal("Error reading from the database.")
	}
}
