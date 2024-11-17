package skiplist

import (
	"bytes"
	"math"
	randv2 "math/rand/v2"
)

// using p = 1/2, i.e. 50% of the nodes are promoted to the next level
// This means that max level is log_2(n), where n is the max number of nodes in the list
// To keep things simple we say that max number of nodes we need to support is 2^16 = 65536
// and log_2(65536) = 16 so that is the max height.

// Probablity of a node occupying level l (zero indexed) is 1/2^l
// So if we're using a psuedo random number generator that uses the full uint32 value range
// we can slice that range in half repreatedly and store the largest value of the lower range
// in an array that has one value for each level. We can then generate a random number
// and go through the array until we find a value that is larger than the random number. The
// index of that value is the level of the new node.

const (
	MaxHeight = 16
	PValue    = 0.5
)

var probabilities [MaxHeight]uint32

func init() {
	levelProb := 1.0 // level 0 has all values

	for level := 0; level < MaxHeight; level++ {
		probabilities[level] = uint32(levelProb * float64(math.MaxUint32))
		levelProb *= PValue
	}
}

func randomHeight() int {
	randVal := randv2.Uint32()

	height := 1
	for height < MaxHeight && randVal <= probabilities[height] {
		height++
	}

	return height
}

type node struct {
	key   []byte
	value []byte
	tower [MaxHeight]*node
}
type SkipList struct {
	head   *node
	height int
}

func NewSkipList() *SkipList {
	return &SkipList{head: &node{}, height: 1}
}

// Insert inserts a new key-value pair into the list. If the key already exists
// the old value is returned along with true. If the key did not exist false is returned.
func (s *SkipList) Insert(key, value []byte) ([]byte, bool) {
	found, rightmostSmaller := s.search(key)
	if found != nil {
		// The key already exists in the list. Update the value.
		oldValue := found.value
		found.value = value
		return oldValue, true
	}

	newNodeHeight := randomHeight()

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
	newNode := &node{key: key, value: value}
	for level := 0; level < newNodeHeight; level++ {
		newNode.tower[level] = rightmostSmaller[level].tower[level]
		rightmostSmaller[level].tower[level] = newNode
	}

	return nil, false
}

// Delete deletes a key-value pair from the list. If the key was found the old value is returned
// along with true. If the key was not found false is returned.
func (s *SkipList) Delete(key []byte) ([]byte, bool) {
	found, rightmostSmaller := s.search(key)
	if found == nil {
		return nil, false
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
func (s *SkipList) Find(key []byte) ([]byte, bool) {
	found, _ := s.search(key)
	if found != nil {
		return found.value, true
	}
	return nil, false
}

func (s *SkipList) search(key []byte) (*node, [MaxHeight]*node) {
	var next *node
	var rightmostSmaller [MaxHeight]*node

	current := s.head
	for level := s.height - 1; level >= 0; level-- {
		// Go through the list at the current level. If we find a nil node we need to drop down a level.
		for next = current.tower[level]; next != nil; next = current.tower[level] {
			if bytes.Compare(key, next.key) <= 0 {
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

	if next != nil && bytes.Equal(key, next.key) {
		return next, rightmostSmaller
	}
	return nil, rightmostSmaller
}
