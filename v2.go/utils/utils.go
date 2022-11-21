// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func CleanupStrings(s string) string {
	return strings.Join(strings.Fields(strings.Replace(s, "\n", "", -1)), " ")
}

func CalculateFileHash(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Can't open file %v: %v", fileName, err))
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", errors.New(fmt.Sprintf("Can't read file %v: %v", fileName, err))
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	return sum, nil
}
