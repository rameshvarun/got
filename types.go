package main

import (
	"crypto/sha256"
	"time"
)

// Hash represents the result of a SHA-256 Hash.
type Hash [32]byte

// Commit represents a commit object.
type CommitObject struct {
	// The author of this commit
	Author string
	// The commit message
	Message string
	// When the commit was created
	Time time.Time
	// The hash of the tree object associated with this commit
	Tree Hash
	// The parents of this commit
	Parents []Hash
}

// Tree represents a tree object.
type TreeObject struct {
	Files map[string]Hash
}

func (this TreeObject) CalculateHash() Hash {
	if len(this.Files) > 0 {

	}

	// Empty directory corresponds to empty hash
	return EMPTY
}

type Difference struct {
	Type     string
	FilePath string
}

// Represents the hash of an empty object
var EMPTY Hash

func init() {
	EMPTY = sha256.Sum256([]byte{})
}
