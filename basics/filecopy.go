package basics

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Copy contents from src to dst file
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// CopyAndHash copies contents from src to dst file and calculates a checksum for it
func CopyAndHash(src, dst string) ([]byte, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return nil, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer source.Close()

	hash := sha256.New()
	sourceWithHash := io.TeeReader(source, hash)

	destination, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer destination.Close()
	_, err = io.Copy(destination, sourceWithHash)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), err
}
