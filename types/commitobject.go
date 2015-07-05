package types

import (
	"bytes"
	"encoding/gob"
	"time"
)

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
func (commit *CommitObject) Serialize() []byte {
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
func DeserializeCommitObject(input []byte) *CommitObject {
	buffer := bytes.NewBuffer(input)
	dec := gob.NewDecoder(buffer)

	var commit CommitObject
	err := dec.Decode(&commit)
	if err != nil {
		panic(err)
	}
	return &commit
}
