package skiplist

import (
	"bytes"
	"cmp"
	"math"
	randv2 "math/rand/v2"
	"testing"

	"github.com/kkdai/basiclist"
)

func TestSkipList_Insert(t *testing.T) {
	s := NewSkipList[[]byte, string](bytes.Compare)

	key := []byte("key1")
	value := "value1"

	// Test inserting a new key
	oldValue, updated := s.Insert(key, value)
	if updated {
		t.Errorf("expected updated to be false, got true")
	}
	if oldValue != "" {
		t.Errorf("expected oldValue to be empty, got %v", oldValue)
	}

	// Test finding the inserted key-value
	foundValue, found := s.Find(key)
	if !found {
		t.Errorf("expected found to be true, got false")
	}
	if foundValue != value {
		t.Errorf("expected oldValue to be empty, got %v", oldValue)
	}

	// Test inserting the same key again
	newValue := "value2"
	oldValue, updated = s.Insert(key, newValue)
	if !updated {
		t.Errorf("expected updated to be true, got false")
	}
	if oldValue != value {
		t.Errorf("expected oldValue to be %v, got %v", value, oldValue)
	}
}

func TestSkipList_Delete(t *testing.T) {
	s := NewSkipList[string, []byte](cmp.Compare[string])

	key := "key1"
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

	// Test finding the deleted key
	_, found := s.Find(key)
	if found {
		t.Errorf("expected found to be false, got true")
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
	s := NewSkipList[int, []byte](cmp.Compare[int])

	key := 1
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

func BenchmarkSkipList_Insert(b *testing.B) {
	s := NewSkipList[int, int](cmp.Compare[int])
	rndVals := randomIntValues(MaxElements)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Insert(rndVals[i%MaxElements], i)
	}
}

func BenchmarkSkipList_Delete(b *testing.B) {
	s := NewSkipList[int, int](cmp.Compare[int])
	rndVals := randomIntValues(MaxElements)
	elemCount := int(math.Min(float64(b.N), MaxElements))

	for i := 0; i < elemCount; i++ {
		s.Insert(rndVals[i], i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Delete(rndVals[i%MaxElements])
	}
}

func BenchmarkSkipList_Find(b *testing.B) {
	s := NewSkipList[int, int](cmp.Compare[int])
	rndVals := randomIntValues(MaxElements)
	elemCount := int(math.Min(float64(b.N), MaxElements))

	for i := 0; i < elemCount; i++ {
		s.Insert(rndVals[i], i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Find(rndVals[i%MaxElements])
	}
}

func BenchmarkLinkedList_Insert(b *testing.B) {
	l := basiclist.NewBasicList()
	rndVals := randomIntValues(MaxElements)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Insert(rndVals[i%MaxElements], i)
	}
}

func BenchmarkLinkedList_Delete(b *testing.B) {
	l := basiclist.NewBasicList()
	rndVals := randomIntValues(MaxElements)
	elemCount := int(math.Min(float64(b.N), MaxElements))

	for i := 0; i < elemCount; i++ {
		l.Insert(rndVals[i], i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Remove(rndVals[i%MaxElements])
	}
}

func BenchmarkLinkedList_Find(b *testing.B) {
	l := basiclist.NewBasicList()
	rndVals := randomIntValues(MaxElements)
	elemCount := int(math.Min(float64(b.N), MaxElements))

	for i := 0; i < elemCount; i++ {
		l.Insert(rndVals[i], i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Search(rndVals[i%MaxElements])
	}
}


func randomIntValues(n int) []int {
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = randv2.Int()
	}
	return values
}

