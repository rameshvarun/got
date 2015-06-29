package main

import (
	"crypto/sha256"
	"io/ioutil"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

func TakeSnapshot(b *bolt.Bucket, filePath string) Hash {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Fatalln("Could not call os.Stat on " + filePath)
	}
	if info.IsDir() {

	} else {
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err.Error())
		}

		hasher := sha256.New()
		hash := Hash{}
		copy(hasher.Sum(file), hash[0:32])

		// Put file into database, addressed by hash
		b.Put(hash[0:32], file)

		return hash
	}
	return EMPTY
}

// Commit implements `got commit`
func Commit(c *cli.Context) {
	db := openDB()
	defer db.Close()

	// Perform operations in a read-only lock
	err := db.View(func(tx *bolt.Tx) error {
		TakeSnapshot(tx.Bucket(OBJECTS), ".")

		return nil
	})
	if err != nil {
		log.Fatal("Error reading from the database.")
	}
}
