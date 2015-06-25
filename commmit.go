package main

import (
	"time"
)

// Hash represents the result of a SHA-256 Hash.
type Hash [32]byte

// Commit represents a commit object.
type Commit struct {
	Author  string    // The author of this commit
	Message string    // The commit message
	Time    time.Time // When the commit was created
	Tree    Hash      // The hash of the tree object associated with this commit
	Parents []Hash    // The parents of this commit
}

// Tree represents a tree object.
type Tree struct {
	Files map[string]Hash
}
