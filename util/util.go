package util

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/boltdb/bolt"
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
func OpenDB() *bolt.DB {
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

// DebugLog prints strings iff the local variable debug is set
func DebugLog(a string) {
	if debug {
		fmt.Println(a)
	}
}
