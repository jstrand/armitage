package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jstrand/armitage/basics"
)

func withRoot(root string, files []string) []string {
	var result []string
	for _, file := range files {
		result = append(result, filepath.Join(root, file))
	}
	return result
}

func findNewFilesFromChecksumFolder(checksumPath string, path string) []string {
	var oldFiles []string
	checksumFiles, _ := basics.FilesUnderPath(checksumPath)
	for _, checksumFile := range checksumFiles {
		checksums := basics.ReadChecksumFile(checksumFile)
		files := basics.FilePaths(checksums)
		oldFiles = append(oldFiles, withRoot(path, files)...)
	}

	currentFiles, _ := basics.FilesUnderPath(path)
	newFiles := basics.Subtract(currentFiles, oldFiles)
	return newFiles
}

func copyFiles(files []string, srcRoot, dstRoot string, checksumWriter io.StringWriter) (success int, failures int) {
	errorCount := 0
	successCount := 0

	for _, filePath := range files {
		pathWithinSourceTree, _ := filepath.Rel(srcRoot, filePath)
		pathToDst := filepath.Join(dstRoot, pathWithinSourceTree)
		dirInDst := filepath.Dir(pathToDst)
		os.MkdirAll(dirInDst, os.ModeDir)
		checksum, err := basics.CopyAndHash(filePath, pathToDst)
		if err != nil {
			log.Println(err)
			errorCount++
		} else {
			successCount++
			checksumWriter.WriteString(basics.FormatShaSum(basics.FileChecksum{Path: pathWithinSourceTree, Checksum: checksum}))
			checksumWriter.WriteString("\n")
		}
	}

	return successCount, errorCount
}

func newChecksumFile(checksumsRoot string) *os.File {
	now := time.Now().Format("2006-01-02T150405")
	newChecksumFilePath := filepath.Join(checksumsRoot, now)
	checksumFile, err := os.OpenFile(newChecksumFilePath, os.O_CREATE, os.ModeAppend)
	if err != nil {
		log.Fatalf("Failed to create new checksum file %s, %s", newChecksumFilePath, err)
	}
	return checksumFile
}

func printUsage() {
	fmt.Println("Armitage")
	fmt.Println("Make backups of new files from given directory.")
	fmt.Println("File copies are put in a directory called <files>.")
	fmt.Println("A folder called <checksums> contains checksums of already copied files.")
	fmt.Println("Changed or previously copied files are ignored (useful for copying photos/video).")
	fmt.Println("Integrity of files can be checked with 'shasum --check ../checksums/*' from the files folder")
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Println("armitage <path-to-copy>")
}

func main() {

	if len(os.Args) != 2 {
		printUsage()
		return
	}

	srcRoot := os.Args[1]
	dstRoot := "files"
	checksumsRoot := "checksums"

	os.Mkdir(dstRoot, os.ModeDir)
	os.Mkdir(checksumsRoot, os.ModeDir)

	newFiles := findNewFilesFromChecksumFolder(checksumsRoot, srcRoot)

	fmt.Printf("%d new files found.", len(newFiles))
	fmt.Println()

	if len(newFiles) == 0 {
		fmt.Println("Done! Nothing to do.")
		return
	}

	checksumFile := newChecksumFile(checksumsRoot)
	defer checksumFile.Close()

	fmt.Printf("Copying from %s to %s\n", srcRoot, dstRoot)
	successCount, errorCount := copyFiles(newFiles, srcRoot, dstRoot, checksumFile)

	if errorCount != 0 {
		fmt.Printf("Warning! %d errors encountered.", errorCount)
		fmt.Println()
	}

	fmt.Printf("Done! Copied %d/%d files.", successCount, len(newFiles))
	fmt.Println()
}
