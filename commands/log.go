package commands

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
	"github.com/rameshvarun/got/util"
)

func Log(c *cli.Context) {
	db := util.OpenDB()
	defer db.Close()

	// Perform operations in a read-only lock
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.INFO)
		heads := b.Get(util.HEADS)

		if heads != nil {
		}

		fmt.Println("There are no commits in this repository...")
		return nil
	})
}
