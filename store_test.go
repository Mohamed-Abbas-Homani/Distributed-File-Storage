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
	key := "Mash"
	expectedFilename := "453f2e740bd7f436e5622609bed9db9edc000bf7"
	expectedPathName := "453f2/e740b/d7f43/6e562/2609b/ed9db/9edc0/00bf7"
	pathKey := CASPathTransformFunc(key)
	if pathKey.Pathname != expectedPathName {
		t.Errorf("have %s want %s", pathKey.Pathname, expectedPathName)
	}
	if pathKey.Filename != expectedFilename {
		t.Errorf("have %s want %s", pathKey.Filename, expectedFilename)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)
	data := []byte("some jpeg bytes")
	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("Test_%d", i)
		if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			t.Error(err)
		}

		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected to NOT have key %s", key)
		}
	}
}
