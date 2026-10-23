package storage

import (
	"sync"
	"testing"
)

func TestPut_NewEntry(t *testing.T) {
	s := &StorageMap{}
	s.Put("key1", StorageValue{DataType: "string", Data: "value1"})

	value, ok := s.Get("key1")
	if !ok {
		t.Fatalf("expected key1 to exist")
	}
	if value.DataType != "string" || value.Data != "value1" {
		t.Fatalf("unexpected value: %+v", value)
	}
}

func TestPut_OverwriteEntry(t *testing.T) {
	s := &StorageMap{}
	s.Put("key1", StorageValue{DataType: "string", Data: "value1"})
	s.Put("key1", StorageValue{DataType: "int", Data: "42"})

	value, ok := s.Get("key1")
	if !ok {
		t.Fatalf("expected key1 to exist")
	}
	if value.DataType != "int" || value.Data != "42" {
		t.Fatalf("unexpected value after overwrite: %+v", value)
	}
}

func TestPut_NilMapInitialization(t *testing.T) {
	var s StorageMap // zero value, data is nil
	s.Put("key1", StorageValue{DataType: "boolean", Data: "true"})

	value, ok := s.Get("key1")
	if !ok {
		t.Fatalf("expected key1 to exist after Put on nil map")
	}
	if value.DataType != "boolean" || value.Data != "true" {
		t.Fatalf("unexpected value: %+v", value)
	}
}

func TestGet_ExistingKey(t *testing.T) {
	s := &StorageMap{}
	s.Put("foo", StorageValue{DataType: "char", Data: "a"})

	value, ok := s.Get("foo")
	if !ok {
		t.Fatalf("expected foo to exist")
	}
	if value.DataType != "char" || value.Data != "a" {
		t.Fatalf("unexpected value: %+v", value)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := &StorageMap{}
	s.Put("foo", StorageValue{DataType: "char", Data: "a"})

	_, ok := s.Get("missing")
	if ok {
		t.Fatalf("expected missing key to not be found")
	}
}

func TestGet_NilMap(t *testing.T) {
	var s StorageMap // zero value, data is nil

	value, ok := s.Get("anything")
	if ok {
		t.Fatalf("expected ok to be false on nil map")
	}
	if value != (StorageValue{}) {
		t.Fatalf("expected zero value, got: %+v", value)
	}
}

func TestClean_ExistingKey(t *testing.T) {
	s := &StorageMap{}
	s.Put("foo", StorageValue{DataType: "string", Data: "bar"})

	s.Clean("foo")

	_, ok := s.Get("foo")
	if ok {
		t.Fatalf("expected foo to be removed after Clean")
	}
}

func TestClean_MissingKey(t *testing.T) {
	s := &StorageMap{}
	// Cleaning a key that was never set should not panic.
	s.Clean("missing")

	_, ok := s.Get("missing")
	if ok {
		t.Fatalf("expected missing key to not exist")
	}
}

func TestClean_NilMap(t *testing.T) {
	var s StorageMap // zero value, data is nil
	// delete on a nil map is a no-op and must not panic.
	s.Clean("foo")
}

func TestConcurrentAccess(t *testing.T) {
	s := &StorageMap{}
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Put("key", StorageValue{DataType: "int", Data: "1"})
			s.Get("key")
		}(i)
	}
	wg.Wait()

	value, ok := s.Get("key")
	if !ok {
		t.Fatalf("expected key to exist after concurrent writes")
	}
	if value.DataType != "int" {
		t.Fatalf("unexpected value: %+v", value)
	}
}
