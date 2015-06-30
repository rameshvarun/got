package types

import (
	"bytes"
	"crypto/sha256"
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
	hash := sha256.Sum256(bytes)
	return hash[0:32]
}

// Represents the hash of an empty object
var EMPTY Hash

func init() {
	EMPTY = CalculateHash([]byte{})
}
