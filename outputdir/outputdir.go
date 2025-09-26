package outputdir

import (
	"os"
	"path/filepath"
)

type OutputDirectory string

func NewOutputDirectory(path string) (OutputDirectory, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	dir := OutputDirectory(absolutePath)
	return dir.Create("")
}

// Creates the relativePath to OutputDirectory in it. Returns a new OutputDirectory of the relative type.
// You can pass "" to relativePath to Create for OutputDirectory.
func (dir OutputDirectory) Create(relativePath string) (OutputDirectory, error) {
	if dir.isPresent(relativePath) {
		return OutputDirectory(dir.relativeToAbsolute(relativePath)), nil
	}

	joinedPath := dir.relativeToAbsolute(relativePath)
	err := os.MkdirAll(joinedPath, 0750)

	return OutputDirectory(joinedPath), err
}

// Returns if relativePath is present in OutputDirectory.
// To check if OutputDirectory is present pass "" to relativePath.
func (dir OutputDirectory) isPresent(relativePath string) bool {
	joinedPath := dir.relativeToAbsolute(relativePath)

	if _, err := os.Stat(joinedPath); !os.IsNotExist(err) {
		return true
	}

	return false
}

// Returns absolute path for a relative path from the OutputDirectory.
func (dir OutputDirectory) relativeToAbsolute(relativePath string) string {
	joinedPath := filepath.Join(dir.String(), relativePath)

	return joinedPath
}

// Writes data to a file in OutputDirectory.
func (dir OutputDirectory) writeFile(filename string, data string) error {
	return os.WriteFile(
		dir.relativeToAbsolute(filename),
		[]byte(data),
		0700,
	)
}

func (dir OutputDirectory) String() string {
	return string(dir)
}
