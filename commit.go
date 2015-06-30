package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

// TakeSnapshot recursivelly snapshots the current directory.
func TakeSnapshot(b *bolt.Bucket, filePath string) Hash {
	// Determine if path is a directory or a file
	info, err := os.Stat(filePath)
	if err != nil {
		log.Fatalln("Could not call os.Stat on " + filePath)
	}

	if info.IsDir() {
		// If the path is a directory, we need to create a tree object
		treeObject := NewTreeObject()

		// Enumerate all of the files in this directory
		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			log.Fatalf("Could not list files in dir %s: %v", filePath, err)
		}
		for _, file := range files {
			if IgnorePath(file.Name()) {
				continue
			}

			// Snapshot the file, and add it to the current tree object
			hash := TakeSnapshot(b, path.Join(filePath, file.Name()))
			treeObject.AddFile(file.Name(), hash)
		}

		// Calculate hash of the directory
		serialized := treeObject.Serialize()
		hash := CalculateHash(serialized)

		// Put directory into database, addressed by it's hash
		b.Put(hash, serialized)
		DebugLog("Put directory \"" + filePath + "\" into database with hash " + hash.String())
		return hash
	}

	// Calculate hash of file contents
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	hash := CalculateHash(file)

	// Put file into database, addressed by it's hash
	b.Put(hash, file)
	DebugLog("Put file \"" + filePath + "\" into database with hash " + hash.String())
	return hash
}

// Commit implements `got commit`
func Commit(c *cli.Context) {
	db := openDB()
	defer db.Close()

	if len(c.String("message")) == 0 {
		log.Fatalln("Must supply a commit message.")
	}

	if len(c.String("author")) == 0 {
		log.Fatalln("Must supply a commit author.")
	}

	// Perform operations in a write lock
	err := db.Update(func(tx *bolt.Tx) error {
		info := tx.Bucket(INFO)       // The bucket holding repositry metadata
		objects := tx.Bucket(OBJECTS) // The bucket holding got objects
		current := info.Get(CURRENT)

		// Create a commit object from the current directory snapshot
		DebugLog("Taking snapshot of repo...")
		hash := TakeSnapshot(tx.Bucket(OBJECTS), ".")
		DebugLog("Repo snapshot has hash " + hash.String() + ".")

		commit := CommitObject{
			Author:  c.String("author"),
			Message: c.String("message"),
			Time:    time.Now(),
			Tree:    hash,
			Parents: make([]Hash, 0),
		}

		// There is a 'current' commit, then it is the parent of the new commit
		if current != nil {
			commit.Parents = append(commit.Parents, current)
		}

		commitBytes := commit.Serialize()
		commitSha := CalculateHash(commitBytes)
		objects.Put(commitSha, commitBytes)
		DebugLog("Created commit of sha " + commitSha.String() + ".")

		// CURRENT now corresponds to the commit that we just created
		info.Put(CURRENT, commitSha)

		// TODO: Update heads

		return nil
	})
	if err != nil {
		log.Fatal("Error reading from the database.")
	}
}
