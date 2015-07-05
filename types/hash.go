package types

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
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
	hash := sha1.Sum(bytes)
	return hash[0:20]
}

// Represents the hash of an empty object
var EMPTY Hash

// SerializeHashes serializes a list of hashes
func SerializeHashes(hashes []Hash) []byte {
	buffer := new(bytes.Buffer)
	e := gob.NewEncoder(buffer)
	err := e.Encode(hashes)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

// DeserializeHashes deserializes a list of hashes
func DeserializeHashes(input []byte) []Hash {
	buffer := bytes.NewBuffer(input)
	dec := gob.NewDecoder(buffer)

	var hashes []Hash
	err := dec.Decode(&hashes)
	if err != nil {
		panic(err)
	}
	return hashes
}

func init() {
	EMPTY = CalculateHash([]byte{})
}
