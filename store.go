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

const defaultRootFolderName = "testFolder"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	//fmt.Println("length of hash string: ", hashStr)
	blocksize := 8
	slicelen := len(hashStr) / blocksize
	//fmt.Println(slicelen)

	paths := make([]string, slicelen)

	for i := 0; i < slicelen; i++ {
		from, to := i*blocksize, ((i + 1) * blocksize)
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Read(key string) (io.Reader, error) {
	fp, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, fp)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	path := s.PathTransformFunc(key)
	return os.Open(path.FullPath())
}

func (s *Store) Has(key string) bool {
	path := s.PathTransformFunc(key)
	_, err := os.Stat(path.FullPath())
	//fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, path.FullPath())

	//_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathkey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathkey.FullPath())
	}()
	return os.RemoveAll(pathkey.FirstPathName())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, path.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	//filename := "somefilename"
	// filenameHash := md5.Sum(buf.Bytes())
	// filename := hex.EncodeToString(filenameHash[:])
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, path.FullPath())
	// fullpath := path.FullPath()
	//fmt.Println(fullPathWithRoot)
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk: %s", n, fullPathWithRoot)
	return nil
}
