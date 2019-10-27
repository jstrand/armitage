package basics

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// FileChecksum stores a file together with its checksum
type FileChecksum struct {
	Path     string
	Checksum []byte
}

// FilePaths takes a number of files with checksums and returns their paths
func FilePaths(checksums []FileChecksum) []string {
	var result []string
	for _, checksum := range checksums {
		result = append(result, checksum.Path)
	}
	return result
}

// FormatShaSum returns a string in the format
// <sha256> *<filepath>
func FormatShaSum(checksum FileChecksum) string {
	return fmt.Sprintf("%x *%s", checksum.Checksum, checksum.Path)
}

func checksum(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// ChecksumFiles takes a number of paths and calculates the checksum for each of them
func ChecksumFiles(paths []string) ([]FileChecksum, error) {
	var result []FileChecksum
	for _, path := range paths {
		checksum, err := checksum(path)
		if err != nil {
			return result, err
		}

		result = append(result, FileChecksum{path, checksum})
	}

	return result, nil
}

func readChecksum(line string) FileChecksum {
	fields := strings.Split(line, "*")
	checksum, _ := hex.DecodeString(strings.TrimSpace(fields[0]))
	path := strings.TrimSpace(fields[1])
	return FileChecksum{Checksum: checksum, Path: path}
}

func readChecksums(lines []string) []FileChecksum {
	var result []FileChecksum
	for _, line := range lines {
		result = append(result, readChecksum(line))
	}
	return result
}

func readLines(path string) ([]string, error) {
	var result []string

	f, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

// ReadChecksumFile deserializes a file with lines in the format:
// <sha256> *<filepath>
// Into a slice with each line represented by a FileChecksum
func ReadChecksumFile(path string) []FileChecksum {
	var result []FileChecksum
	lines, _ := readLines(path)
	for _, line := range lines {
		result = append(result, readChecksum(line))
	}
	return result
}
