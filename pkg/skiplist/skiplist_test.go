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
func TestSkipList_InsertAndRemoveRandomElements(t *testing.T) {
	s := NewSkipList[int, int](cmp.Compare[int])
	rndVals := randomIntValues(100)

	// Insert 100 random elements
	for i := 0; i < len(rndVals); i++ {
		s.Insert(rndVals[i], i)
	}

	// Remove the elements one by one and verify the rest are still there
	for i := 0; i < len(rndVals); i++ {
		keyToRemove := rndVals[i]
		s.Delete(keyToRemove)

		// Verify the removed element is not found
		_, found := s.Find(keyToRemove)
		if found {
			t.Errorf("expected key %v to be removed, but it was found", keyToRemove)
		}

		// Verify the rest of the elements are still there
		for j := i + 1; j < len(rndVals); j++ {
			keyToFind := rndVals[j]
			value, found := s.Find(keyToFind)
			if !found {
				t.Errorf("expected key %v to be found, but it was not", keyToFind)
			}
			if value != j {
				t.Errorf("expected value %v for key %v, got %v", j, keyToFind, value)
			}
		}
	}
}

func BenchmarkSkipList_Insert(b *testing.B) {
	rndVals := randomIntValuesB(b)
	var s *SkipList[int, int]

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rndValIx := i % MaxElements
		if rndValIx == 0 {
			// Starting from the beginning of the random values list, create a new skip list
			// to avoid just inserting keys that are already in the list.
			s = NewSkipList[int, int](cmp.Compare[int])
		}
		b.StartTimer()
		s.Insert(rndVals[rndValIx], i)
	}
}

func BenchmarkSkipList_Delete(b *testing.B) {
	rndVals := randomIntValuesB(b)
	var s *SkipList[int, int]

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rndValIx := i % MaxElements
		if rndValIx == 0 {
			// Starting from the beginning of the random values list, create a new skip list
			// and populate it with random values so were not deleting from an empty list.
			s = NewSkipList[int, int](cmp.Compare[int])
			for i := 0; i < len(rndVals); i++ {
				s.Insert(rndVals[i], i)
			}
		}
		b.StartTimer()
		s.Delete(rndVals[rndValIx])
	}
}

func BenchmarkSkipList_Find(b *testing.B) {
	s := NewSkipList[int, int](cmp.Compare[int])
	rndVals := randomIntValuesB(b)

	for i := 0; i < len(rndVals); i++ {
		s.Insert(rndVals[i], i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Find(rndVals[i%MaxElements])
	}
}

func BenchmarkLinkedList_Insert(b *testing.B) {
	var l basiclist.BasicList
	rndVals := randomIntValuesB(b)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rndValIx := i % MaxElements
		if rndValIx == 0 {
			// Starting from the beginning of the random values list, create a new linked list
			// to avoid just inserting keys that are already in the list.
			l = *basiclist.NewBasicList()
		}
		b.StartTimer()
		l.Insert(rndVals[rndValIx], i)
	}
}

func BenchmarkLinkedList_Delete(b *testing.B) {
	var l basiclist.BasicList
	rndVals := randomIntValuesB(b)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rndValIx := i % MaxElements
		if rndValIx == 0 {
			// Starting from the beginning of the random values list, create a new linked list
			// and populate it with random values so were not deleting from an empty list.
			l = *basiclist.NewBasicList()
			for i := 0; i < len(rndVals); i++ {
				l.Insert(rndVals[i], i)
			}
		}
		b.StartTimer()
		l.Remove(rndVals[rndValIx])
	}
}

func BenchmarkLinkedList_Find(b *testing.B) {
	l := basiclist.NewBasicList()
	rndVals := randomIntValuesB(b)
	elemCount := int(math.Min(float64(b.N), MaxElements))

	for i := 0; i < elemCount; i++ {
		l.Insert(rndVals[i], i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Search(rndVals[i%MaxElements])
	}
}

func randomIntValuesB(b *testing.B) []int {
	return randomIntValues(int(math.Min(float64(b.N), MaxElements)))
}

func randomIntValues(elemCount int) []int {
	values := make([]int, elemCount)
	for i := 0; i < elemCount; i++ {
		values[i] = randv2.Int()
	}
	return values
}
