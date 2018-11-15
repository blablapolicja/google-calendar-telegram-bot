package util

import (
	"encoding/base64"
	"math/rand"
)

// GetRandomStateString returns random string
func GetRandomStateString() string {
	b := make([]byte, 32)

	rand.Read(b)

	return base64.StdEncoding.EncodeToString(b)
}
