package skiplist

import (
	"math"
	randv2 "math/rand/v2"
)

type node[K any, V any] struct {
	key   K
	value V
	tower []*node[K, V]
}
type SkipList[K any, V any] struct {
	head          *node[K, V]
	height        int
	maxHeight     int
	probabilities []uint32
	comparator    func(a, b K) int
}

func NewDefaultSkipList[K any, V any](keyComparator func(a, b K) int) *SkipList[K, V] {
	return NewSkipList[K, V](keyComparator, 65536, 0.5)
}

// NewSkipList creates a new SkipList with the given key comparator function.
// The key comparator function should return a negative value if a < b, 0 if a == b
// and a positive value if a > b.
func NewSkipList[K any, V any](keyComparator func(a, b K) int, performantCapacity int, pValue float64) *SkipList[K, V] {
	if performantCapacity < 1 {
		panic("performantCapacity must be at least 1")
	}
	if pValue <= 0.0 || pValue >= 1.0 {
		panic("pValue must be in the range (0, 1)")
	}
	maxHeight := int(math.Ceil(logBaseX(1.0/pValue, float64(performantCapacity))))

	answer := &SkipList[K, V]{
		head:          &node[K, V]{tower: make([]*node[K, V], maxHeight)},
		height:        1,
		maxHeight:     maxHeight,
		probabilities: make([]uint32, maxHeight),
		comparator:    keyComparator,
	}

	// Probablity of a node occupying level l (zero indexed) is pValue^l
	// So if we're using a psuedo random number generator that uses the full uint32 value range
	// we can slice that range in half repreatedly and store the largest value of the lower range
	// in an array that has one value for each level. We can then generate a random number
	// and go through the array until we find a value that is larger than the random number. The
	// index of that value is the level of the new node.
	levelProb := 1.0 // level 0 has all values
	for level := 0; level < maxHeight; level++ {
		answer.probabilities[level] = uint32(levelProb * float64(math.MaxUint32))
		levelProb *= pValue
	}
	return answer
}

func logBaseX(base, x float64) float64 {
	return math.Log(x) / math.Log(base)
}

func (s *SkipList[K, V]) randomHeight() int {
	randVal := randv2.Uint32()

	height := 1
	for height < s.maxHeight && randVal <= s.probabilities[height] {
		height++
	}

	return height
}

// Insert inserts a new key-value pair into the list. If the key already exists
// the old value is returned along with true. If the key did not exist false is returned.
func (s *SkipList[K, V]) Insert(key K, value V) (V, bool) {
	found, rightmostSmaller := s.search(key)
	if found != nil {
		// The key already exists in the list. Update the value.
		oldValue := found.value
		found.value = value
		return oldValue, true
	}

	newNodeHeight := s.randomHeight()

	if newNodeHeight > s.height {
		// The new node is taller than the current list. This means that head will be the previous
		// node for the new levels.
		for newLevel := s.height; newLevel < newNodeHeight; newLevel++ {
			rightmostSmaller[newLevel] = s.head
		}
		s.height = newNodeHeight
	}

	// Insert a new node and point rightmostSmaller nodes at each level to the new node (up to
	// the height of the new node).
	newNode := &node[K, V]{
		key:   key,
		value: value,
		tower: make([]*node[K, V], s.maxHeight),
	}
	for level := 0; level < newNodeHeight; level++ {
		newNode.tower[level] = rightmostSmaller[level].tower[level]
		rightmostSmaller[level].tower[level] = newNode
	}

	return *new(V), false
}

// Delete deletes a key-value pair from the list. If the key was found the old value is returned
// along with true. If the key was not found false is returned.
func (s *SkipList[K, V]) Delete(key K) (V, bool) {
	found, rightmostSmaller := s.search(key)
	if found == nil {
		return *new(V), false
	}

	// Start from level 0 and see if the rightmost node with a smaller key at this level
	// points directly to the node we're deleting. If it does, update the pointer to point
	// to the next node. If it does not, we're done since it means we have reached the height
	// of the node we're deleting.
	for level := 0; level < s.height; level++ {
		if rightmostSmaller[level].tower[level] != found {
			break
		}
		rightmostSmaller[level].tower[level] = found.tower[level]
	}

	// Update the height of the list if the node we're deleting is the highest node in the list.
	for s.height > 1 && s.head.tower[s.height-1] == nil {
		s.height--
	}

	return found.value, true
}

// Find finds a value in the list given its key. If the key is found the value is returned
// along with true, otherwise false is returned.
func (s *SkipList[K, V]) Find(key K) (V, bool) {
	found, _ := s.search(key)
	if found != nil {
		return found.value, true
	}
	return *new(V), false
}

func (s *SkipList[K, V]) search(key K) (*node[K, V], []*node[K, V]) {
	var next *node[K, V]
	rightmostSmaller := make([]*node[K, V], s.maxHeight)

	current := s.head
	for level := s.height - 1; level >= 0; level-- {
		// Go through the list at the current level. If we find a nil node we need to drop down a level.
		for next = current.tower[level]; next != nil; next = current.tower[level] {
			if s.comparator(key, next.key) <= 0 {
				// This means that we have found the node that is either the node we're looking for
				// or a node with a key that is larger than the key we're looking for. Time to go down a level.

				// Even if we found the node we're looking for we need to go all levels to populate rightmostSmaller,
				// which is needed for instertion and deletion. An optimized version could have a separete search
				// implementation for Find that does not keep track of "previous per level" and only goes down to
				// the level where the node is found.
				break
			}
			// The next node at this level has a key that is smaller than the key we're looking for. Move to this node.
			current = next
		}
		// We now know that the current node is the righmost node at the current level that has a smaller key than the
		// one we're looking for.
		rightmostSmaller[level] = current
	}

	if next != nil && s.comparator(key, next.key) == 0 {
		return next, rightmostSmaller
	}
	return nil, rightmostSmaller
}
