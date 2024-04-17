package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "myjourney"
	pathKey := CASPathTransformFunc(key)
	expectedFilename := "1fa3bcce10956c7d4e2ebd7cae42f90657bc4ae5"
	expectedPathName := "1fa3bcce/10956c7d/4e2ebd7c/ae42f906/57bc4ae5"

	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s want %s", pathKey.PathName, expectedPathName)
	}

	if pathKey.FileName != expectedFilename {
		t.Errorf("have %s want %s", pathKey.FileName, expectedFilename)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	data := []byte("some random bytes")
	if err := s.writeStream("randombytes", bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	r, err := s.Read("randombytes")
	if err != nil {
		t.Error(err)
	}
	b, _ := io.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}

	if err := s.Delete("randombytes"); err != nil {
		t.Error(err)
	}
}
