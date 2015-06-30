package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"sort"
	"time"
)

// Hash represents the result of a SHA-256 Hash.
type Hash []byte

func (hash Hash) String() string {
	return hex.EncodeToString(hash)
}

// Equal tests Equality between two hashes
func (hash Hash) Equal(other Hash) bool {
	return bytes.Equal(hash, other)
}

// CalculateHash calculates the hash of a byte slice
func CalculateHash(bytes []byte) Hash {
	hash := sha256.Sum256(bytes)
	return hash[0:32]
}

// CommitObject represents a commit, with its associated metadata and
// directory snapshot.
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

// Serialize creates a representation of the current CommitObject.
func (commit CommitObject) Serialize() []byte {
	buffer := new(bytes.Buffer)
	e := gob.NewEncoder(buffer)
	err := e.Encode(commit)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

// DeserializeCommitObject returns a commit object deserialized from the given
// byte stream
func DeserializeCommitObject(input []byte) CommitObject {
	buffer := bytes.NewBuffer(input)
	dec := gob.NewDecoder(buffer)

	var commit CommitObject
	err := dec.Decode(&commit)
	if err != nil {
		panic(err)
	}
	return commit
}

// DirectoryEntry represents one entry in a TreeObject
type DirectoryEntry struct {
	Name string
	Hash Hash
}

// TreeObject represents a directory.
type TreeObject struct {
	Files []DirectoryEntry
}

// ByFileName implements sort.Interface for []Person based on
// the Age field.
type ByFileName []DirectoryEntry

func (a ByFileName) Len() int           { return len(a) }
func (a ByFileName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFileName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// NewTreeObject creates a new TreeObject
func NewTreeObject() TreeObject {
	return TreeObject{
		Files: make([]DirectoryEntry, 0),
	}
}

// AddFile adds a file to the given TreeObject
func (tree TreeObject) AddFile(name string, hash Hash) {
	tree.Files = append(tree.Files, DirectoryEntry{
		Name: name,
		Hash: hash,
	})
	sort.Sort(ByFileName(tree.Files))
}

// HasFile returns true if the tree has a file of the given name
func (tree TreeObject) HasFile(name string) bool {
	for _, entry := range tree.Files {
		if entry.Name == name {
			return true
		}
	}
	return false
}

// Serialize creates a representation of the current TreeObject.
func (tree TreeObject) Serialize() []byte {
	if len(tree.Files) > 0 {
		buffer := new(bytes.Buffer)
		e := gob.NewEncoder(buffer)

		err := e.Encode(tree)
		if err != nil {
			panic(err)
		}
		return buffer.Bytes()
	}

	// Empty directory corresponds to empty bytes
	return []byte{}
}

// DeserializeTreeObject returns a commit object deserialized from the given
// byte stream
func DeserializeTreeObject(input []byte) TreeObject {
	buffer := bytes.NewBuffer(input)
	dec := gob.NewDecoder(buffer)

	var tree TreeObject
	err := dec.Decode(&tree)
	if err != nil {
		panic(err)
	}
	return tree
}

// Difference represents a modification between two directory trees.
type Difference struct {
	Type     string
	FilePath string
}

// Represents the hash of an empty object
var EMPTY Hash

func init() {
	EMPTY = CalculateHash([]byte{})
}
