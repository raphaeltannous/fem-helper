package cache

import (
	"errors"
	"os"
	"path/filepath"
)

var cache CacheDir
var cacheError error

type CacheDir string

func NewCache() (CacheDir, error) {
	if cache != "" {
		return cache, cacheError
	}

	var cdir string
	cdir, cacheError = os.UserCacheDir()
	if cacheError != nil {
		return "", cacheError
	}
	cache = CacheDir(filepath.Join(cdir, "fem-helper"))

	_, err := os.Stat(string(cache))
	if errors.Is(err, os.ErrNotExist) {
		cacheError = os.MkdirAll(string(cache), 0700)
	}

	return cache, cacheError
}

// Save data to filename in cache and returns the number of bytes writing.
func (cache *CacheDir) Save(filename string, data []byte) (int, error) {
	cacheFile, err := os.Create(cache.GetAbsolutePath(filename))
	if err != nil {
		return 0, err
	}
	defer cacheFile.Close()

	return cacheFile.Write(data)
}

// Returns the content of a given filename in cache if it exists.
func (cache *CacheDir) Read(filename string) ([]byte, error) {
	cacheFile := cache.GetAbsolutePath(filename)
	return os.ReadFile(cacheFile)
}

// Returns absolute path given a relativePath to the cache.
func (cache *CacheDir) GetAbsolutePath(relativePath string) string {
	return filepath.Join(string(*cache), relativePath)
}

func (cache CacheDir) String() string {
	return string(cache)
}
