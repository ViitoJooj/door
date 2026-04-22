package initproject

import (
	"crypto/rand"
	"encoding/hex"
)

func randomHex(size int) string {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
