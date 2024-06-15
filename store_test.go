package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "Mash"
	expectedFilename := "453f2e740bd7f436e5622609bed9db9edc000bf7"
	expectedPathName := "453f2/e740b/d7f43/6e562/2609b/ed9db/9edc0/00bf7"
	pathkey := CASPathTransformFunc(key)
	if pathkey.Pathname != expectedPathName {
		t.Errorf("have %s want %s", pathkey.Pathname, expectedPathName)
	}
	if pathkey.Filename != expectedFilename {
		t.Errorf("have %s want %s", pathkey.Filename, expectedFilename)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOtps{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	data := []byte("some jpeg bytes")
	key := "myspecialimage"
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
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

    s.Delete(key)

}

func TestStoreDeleteKey ( t *testing.T) {
	opts := StoreOtps{
		PathTransformFunc: CASPathTransformFunc,
	}
    s := NewStore(opts)
	data := []byte("some jpeg bytes")
	key := "myspecialimage"
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

    if err := s.Delete(key); err != nil {
        t.Error(err)
    }
}