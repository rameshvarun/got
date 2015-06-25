package main

import (
	"time"
)

// Hash represents the result of a SHA-256 Hash.
type Hash [32]byte

// Commit represents a commit object.
type Commit struct {
	author      string    // The author of this commit
	description string    // The commit description
	time        time.Time // When the commit was created
	tree        Hash      // The hash of the tree object associated with this commit
	parents     []Hash    // The parents of this commit
}

// Tree represents a tree object.
type Tree struct {
	files map[string]Hash
}
