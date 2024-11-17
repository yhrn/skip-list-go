package skiplist

import (
	"bytes"
	"testing"
)

func TestSkipList_Insert(t *testing.T) {
	s := NewSkipList()

	key := []byte("key1")
	value := []byte("value1")

	// Test inserting a new key
	oldValue, updated := s.Insert(key, value)
	if updated {
		t.Errorf("expected updated to be false, got true")
	}
	if oldValue != nil {
		t.Errorf("expected oldValue to be nil, got %v", oldValue)
	}

	// Test inserting the same key again
	newValue := []byte("value2")
	oldValue, updated = s.Insert(key, newValue)
	if !updated {
		t.Errorf("expected updated to be true, got false")
	}
	if !bytes.Equal(oldValue, value) {
		t.Errorf("expected oldValue to be %v, got %v", value, oldValue)
	}
}

func TestSkipList_Delete(t *testing.T) {
	s := NewSkipList()

	key := []byte("key1")
	value := []byte("value1")

	// Test deleting a non-existent key
	oldValue, deleted := s.Delete(key)
	if deleted {
		t.Errorf("expected deleted to be false, got true")
	}
	if oldValue != nil {
		t.Errorf("expected oldValue to be nil, got %v", oldValue)
	}

	// Insert a key and then delete it
	s.Insert(key, value)
	oldValue, deleted = s.Delete(key)
	if !deleted {
		t.Errorf("expected deleted to be true, got false")
	}
	if !bytes.Equal(oldValue, value) {
		t.Errorf("expected oldValue to be %v, got %v", value, oldValue)
	}

	// Test deleting the same key again
	oldValue, deleted = s.Delete(key)
	if deleted {
		t.Errorf("expected deleted to be false, got true")
	}
	if oldValue != nil {
		t.Errorf("expected oldValue to be nil, got %v", oldValue)
	}
}

func TestSkipList_Find(t *testing.T) {
	s := NewSkipList()

	key := []byte("key1")
	value := []byte("value1")

	// Test finding a non-existent key
	foundValue, found := s.Find(key)
	if found {
		t.Errorf("expected found to be false, got true")
	}
	if foundValue != nil {
		t.Errorf("expected foundValue to be nil, got %v", foundValue)
	}

	// Insert a key and then find it
	s.Insert(key, value)
	foundValue, found = s.Find(key)
	if !found {
		t.Errorf("expected found to be true, got false")
	}
	if !bytes.Equal(foundValue, value) {
		t.Errorf("expected foundValue to be %v, got %v", value, foundValue)
	}
}
