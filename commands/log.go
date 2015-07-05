package commands

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
	"github.com/rameshvarun/got/types"
	"github.com/rameshvarun/got/util"
)

func Log(c *cli.Context) {
	db := util.OpenDB()
	defer db.Close()

	yellow := ansi.ColorFunc("yellow+h:black")

	// Perform operations in a read-only lock
	db.View(func(tx *bolt.Tx) error {
		info := tx.Bucket(util.INFO)
		objects := tx.Bucket(util.OBJECTS)
		headsBytes := info.Get(util.HEADS)

		if headsBytes != nil {
			queue := types.DeserializeHashes(headsBytes)

			// TODO: Keep a visited set, so we don't repeat commits

			for len(queue) > 0 {
				i := GetNewestCommit(objects, queue)

				// Get commit and remove it from the priority queue
				commitSha := queue[i]
				commit := types.DeserializeCommitObject(objects.Get(commitSha))
				queue = append(queue[:i], queue[i+1:]...)

				fmt.Printf(yellow("commit %s\n"), commitSha)
				fmt.Printf("Message: %s\n", commit.Message)
				fmt.Printf("Author: %s\n", commit.Author)
				fmt.Printf("Date: %s\n", commit.Time)

				if len(commit.Parents) > 0 {
					fmt.Printf("Parents: %s\n", commit.Parents)
				}

				fmt.Printf("Tree: %s\n", commit.Tree)
				fmt.Println()

				// Append parents of this commit to the queue
				queue = append(queue, commit.Parents...)
			}

			return nil
		}

		fmt.Println("There are no commits in this repository...")
		return nil
	})
}

type newestCommit struct {
	Index int
	Time  time.Time
}

func GetNewestCommit(objects *bolt.Bucket, hashes []types.Hash) int {
	var newest *newestCommit = nil

	for i, hash := range hashes {
		commit := types.DeserializeCommitObject(objects.Get(hash))
		if newest == nil || commit.Time.After(newest.Time) {
			newest = new(newestCommit)
			newest.Index = i
			newest.Time = commit.Time
		}
	}
	return newest.Index
}
