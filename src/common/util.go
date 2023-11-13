package common

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString takes a string as input and returns its SHA-256 hash as a hex string.
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return hex.EncodeToString(bs)
}
