package util

import (
	"encoding/base64"
	"math/rand"
)

// GenerateStateString returns random string
func GenerateStateString() string {
	b := make([]byte, 32)

	rand.Read(b)

	return base64.StdEncoding.EncodeToString(b)
}
