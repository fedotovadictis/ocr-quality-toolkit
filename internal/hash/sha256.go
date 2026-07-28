package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func FileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}
	result := hex.EncodeToString(hasher.Sum(nil))
	return result, nil
}
