package env

import (
	"os"
)

func DecodeSecret(key []byte) ([]byte, error) {
	return []byte(os.Getenv(string(key))), nil
}
