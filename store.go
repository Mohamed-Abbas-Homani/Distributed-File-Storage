package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "MashNetwork"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type PathKey struct {
	Pathname string
	Filename string
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

func (p *PathKey) FirstPathName() string {
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}

	return paths[0]
}

type StoreOtps struct {
	// root is the folder name of the root, containing all
	// the files of the system.
	Root              string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOtps
}

func NewStore(opts StoreOtps) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}

	return &Store{
		StoreOtps: opts,
	}
}

func (s *Store) Clear() error {

	defer func() {
		log.Printf("All deleted\n")
	}()
	return os.RemoveAll(s.Root)
}

func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformFunc(key)
    fullPath := fmt.Sprintf("%s/%s", s.Root, pathkey.FullPath())
	_, err := os.Stat(fullPath)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathkey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathkey.Filename)
	}()
	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathkey := s.PathTransformFunc(key)
	fullPath := fmt.Sprintf("%s/%s", s.Root, pathkey.FullPath())
	return os.Open(fullPath)
}

func (s *Store) Write (key string, r io.Reader) (int64, error) {
    return s.writeStream(key, r)
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathkey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return 0, err
	}

	fullPath := fmt.Sprintf("%s/%s", s.Root, pathkey.FullPath())
	f, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}
