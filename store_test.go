package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

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
	s := newStore()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("file_%d", i)

		data := []byte("some random bytes")
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("Got: key %s is not present. Expected: Key %s is present ", key, key)
		}

		r, err := s.Read("randombytes")
		if err != nil {
			t.Error(err)
		}
		b, _ := io.ReadAll(r)

		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}
		if ok := s.Has(key); ok {
			t.Errorf("Got: key %s is present. Expected: Key %s not present ", key, key)
		}
	}
}
