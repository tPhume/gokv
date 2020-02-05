package btree

import (
	"errors"
	"github.com/tPhume/gokv/store"
)

// Package contains in memory implementation of btree

var (
	KeyDoesNotExist = errors.New("key does not exist")
)

// holds the key and value pair
type item struct {
	key   string
	value store.Value
}

func (i *item) getKey() string {
	return i.key
}

func (i *item) getValue() store.Value {
	return i.value
}

// node holds an slice of items and slice of children nodes
type node struct {
	items     []*item
	node      []*node
	minDegree int
	currKey   int
	leaf      bool
}

func newNode(minDegree int, leaf bool) *node {
	return &node{
		items:     make([]*item, 2*minDegree-1),
		node:      make([]*node, 2*minDegree),
		minDegree: minDegree,
		currKey:   0,
		leaf:      leaf,
	}
}

// can only insert if there is space
func (n *node) insert(it *item) error {
	index := n.currKey - 1
	key := it.getKey()

	if n.leaf {
		// if node is leaf, just insert
		for i := index; i >= 0; i-- {
			index = i
			if n.items[i].getKey() > key {
				n.items[i+1] = n.items[i]
			} else {
				break
			}
		}

		n.items[index+1] = copyItem(it)
		n.currKey = n.currKey + 1
	} else {
		for i := index; i >= 0; i-- {
			index = i
			if n.items[i].getKey() < key {
				break
			}
		}

		if n.node[index+1].currKey == 2*n.minDegree-1 {
			err := n.splitChild(index+1, n.node[index+1])
			if err != nil {
				return err
			}

			if key > n.items[index+1].getKey() {
				index++
			}

			err = n.node[index+1].insert(it)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// utility function to split child, can only split if it is full
func (n *node) splitChild(index int, child *node) error {
	minDegree := n.minDegree

	// create new right node
	biggerNode := newNode(child.minDegree, child.leaf)
	biggerNode.currKey = minDegree - 1

	// move items to right node
	for i := 0; i < minDegree-1; i++ {
		biggerNode.items[i] = child.items[i+minDegree]
		child.items[i+minDegree] = nil
	}

	// copy children as well, if this node is root
	if child.leaf {
		for i := 0; i < minDegree; i++ {
			biggerNode.node[i] = child.node[i+minDegree]
			child.node[i+minDegree] = nil
		}
	}

	child.currKey = minDegree - 1

	// create space and add new node
	for i := n.currKey; i >= index+1; i-- {
		n.node[i+1] = n.node[i]
	}

	n.node[index+1] = biggerNode

	// create space and add new item
	for i := n.currKey - 1; i >= index; i-- {
		n.items[i+1] = n.items[i]
	}

	n.items[index] = child.items[minDegree-1]
	child.items[minDegree-1] = nil

	n.currKey++

	return nil
}

func (n *node) update(key string, value store.Value) error {
	pos := n.findKey(key)
	if pos == -1 {
		return KeyDoesNotExist
	}

	if pos == n.currKey {
		if n.leaf {
			return KeyDoesNotExist
		}

		return n.node[pos].update(key, value)
	}

	if n.items[pos].getKey() == key {
		n.items[pos] = &item{
			key:   key,
			value: copyValue(value),
		}

		return nil
	}

	if n.leaf {
		return KeyDoesNotExist
	}

	return n.node[pos].update(key, value)
}

func (n *node) search(key string) store.Value {
	pos := n.findKey(key)
	if pos == -1 {
		return nil
	}

	if pos == n.currKey {
		if n.leaf {
			return nil
		}

		return n.node[pos].search(key)
	}

	if n.items[pos].getKey() == key {
		return copyValue(n.items[pos].getValue())
	}

	if n.leaf {
		return nil
	}

	return n.node[pos].search(key)
}

func (n *node) remove(key string) error {
	index := n.findKey(key)
	if index == -1 {
		return KeyDoesNotExist
	}

	// key to delete is in current node
	if index < n.currKey && n.items[index].getKey() == key {
		if n.leaf {
			n.removeFromLeaf(index)
		} else {
			err := n.removeFromNonLeaf(index)
			if err != nil {
				return err
			}
		}
	} else {
		if n.leaf {
			return KeyDoesNotExist
		}

		flag := false
		if index == n.currKey {
			flag = true
		}

		if n.node[index].currKey < n.minDegree {
			n.fill(index)
		}

		if flag && index > n.currKey {
			err := n.node[index-1].remove(key)
			if err != nil {
				return err
			}
		} else {
			err := n.node[index].remove(key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *node) removeFromLeaf(index int) {
	for i := index + 1; i < n.currKey; i++ {
		n.items[i-1] = n.items[i]
	}

	n.items[n.currKey-1] = nil
	n.currKey--
}

func (n *node) removeFromNonLeaf(index int) error {
	minDegree := n.minDegree
	it := n.items[index]

	// if child that precedes "it" has at least minDegree keys
	// find pred and and replace "it" and recursively delete pred
	if n.node[index].currKey >= minDegree {
		pred := n.getPred(index)
		n.items[index] = pred
		err := n.node[index].remove(pred.getKey())
		if err != nil {
			return err
		}
	} else if n.node[index+1].currKey >= minDegree {
		// if node[index] has less than minDegree items
		// check if next node has at least minDegree items
		// find successor of and replace item with successor

		succ := n.getSucc(index)
		n.items[index] = succ
		err := n.node[index+1].remove(succ.getKey())
		if err != nil {
			return err
		}
	} else {
		// if both node[index] and node[index+1]
		// has less than minDegree, merge them
		n.merge(index)
		err := n.node[index].remove(it.getKey())
		if err != nil {
			return err
		}
	}

	return nil
}

// utility function that finds index of item greater than or equal to the key
func (n *node) findKey(key string) int {
	if n.isEmpty() {
		return -1
	}

	pos := 0
	for i := 0; i < n.currKey; i++ {
		if n.items[i].getKey() >= key {
			break
		}

		pos++
	}

	return pos
}

// utility function to get last predecessor of node at index
func (n *node) getPred(index int) *item {
	cur := n.node[index]

	for true {
		if cur.leaf {
			break
		}

		cur = cur.node[cur.currKey]
	}

	return cur.items[cur.currKey-1]
}

// utility function to get successor of a key
func (n *node) getSucc(index int) *item {
	cur := n.node[index+1]

	for true {
		if cur.leaf {
			break
		}

		cur = cur.node[0]
	}

	return cur.items[0]
}

// utility function to fill up child node if child has less than minDegree - 1 keys
func (n *node) fill(index int) {
	if index != 0 && n.node[index-1].currKey > n.minDegree {
		n.borrowFromPrev(index)
	} else if index != n.currKey && n.node[index+1].currKey >= n.minDegree {
		n.borrowFromNext(index)
	} else {
		if index != n.currKey {
			n.merge(index)
		} else {
			n.merge(index - 1)
		}
	}
}

// utility function to borrow key/item from previous node in array
func (n *node) borrowFromPrev(index int) {
	child := n.node[index]
	sibling := n.node[index-1]

	// moving all items in child forward one step
	for i := child.currKey - 1; i >= 0; i-- {
		child.items[i+1] = child.items[i]
	}

	// move child pointers one step ahead
	if !child.leaf {
		for i := child.currKey; i >= 0; i-- {
			child.node[i+1] = child.node[i]
		}
	}

	// add item of index to first item in child
	child.items[0] = n.items[index-1]

	// move last node of sibling to child node
	if !child.leaf {
		child.node[0] = sibling.node[sibling.currKey]
		sibling.node[sibling.currKey] = nil
	}

	n.items[index-1] = sibling.items[sibling.currKey-1]
	sibling.items[sibling.currKey-1] = nil

	child.currKey++
	sibling.currKey--
}

// utility function to borrow key/item from next node in array
func (n *node) borrowFromNext(index int) {
	child := n.node[index]
	sibling := n.node[index+1]

	// insert item[index] from current node as the last item in child
	child.items[child.currKey] = n.items[index]

	// add sibling first node as last node in child
	if !child.leaf {
		child.node[child.currKey+1] = sibling.node[0]
	}

	// add first item from sibling to current node
	n.items[index] = sibling.items[0]

	// moving item in sibling one step back
	for i := 1; i < sibling.currKey; i++ {
		sibling.items[i-1] = sibling.items[i]
	}

	// moving sibling child pointers one step back
	if !sibling.leaf {
		for i := 1; i < sibling.currKey+1; i++ {
			sibling.node[i-1] = sibling.node[i]
		}
	}

	// updating count of the node
	child.currKey++
	sibling.currKey--
}

// merge node at index and index+1
func (n *node) merge(index int) {
	minDegree := n.minDegree

	child := n.node[index]
	sibling := n.node[index+1]

	// add item from current node to child (the middle item)
	child.items[minDegree-1] = n.items[index]

	// copy items from sibling to child
	for i := 0; i < n.currKey; i++ {
		child.items[child.currKey+i] = sibling.items[i]
	}

	// copy child nodes from sibling to child
	if !child.leaf {
		for i := 0; i <= n.currKey; i++ {
			child.node[child.currKey+i] = sibling.node[i]
		}
	}

	// move items to fill gap in current node
	for i := index + 1; i < n.currKey; i++ {
		n.items[i-1] = n.items[i]
	}

	n.items[n.currKey-1] = nil

	// move child nodes to fill in gap
	for i := index + 2; i < n.currKey+1; i++ {
		n.node[i-1] = n.node[i]
	}

	child.currKey = child.currKey + sibling.currKey + 1
	n.currKey--
}

// utility function to check if node is empty
func (n *node) isEmpty() bool {
	return n.currKey == 0
}

// utility function to check if node is full
func (n *node) isFull() bool {
	return n.currKey == 2*n.minDegree-1
}

// btree encapsulates node type which does most of the work
// it also implements the store interface
// the api package will interact with the btree instead of the node type directly
type Btree struct {
	root      *node
	minDegree int
}

func NewBtree(minDegree int) *Btree {
	return &Btree{
		root:      newNode(minDegree, true),
		minDegree: minDegree,
	}
}

func (b *Btree) Insert(key string, value store.Value) error {
	if b.root.isFull() {
		newRoot := newNode(b.minDegree, false)
		newRoot.node[0] = b.root

		err := newRoot.splitChild(0, b.root)
		if err != nil {
			return err
		}

		if key < newRoot.items[0].getKey() {
			err := newRoot.node[0].insert(&item{
				key:   key,
				value: value,
			})

			if err != nil {
				return err
			}
		} else {
			err := newRoot.node[1].insert(&item{
				key:   key,
				value: value,
			})

			if err != nil {
				return err
			}
		}

		b.root = newRoot
	} else {
		err := b.root.insert(&item{
			key:   key,
			value: value,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Btree) Update(key string, value store.Value) error {
	return b.root.update(key, value)
}

func (b *Btree) Search(key string) store.Value {
	return b.root.search(key)
}

func (b *Btree) Remove(key string) error {
	err := b.root.remove(key)
	if err != nil {
		return err
	}

	if b.root.currKey == 0 {
		if !b.root.leaf {
			b.root = b.root.node[0]
		}
	}

	return nil
}

// utility functions
func copyValue(v store.Value) store.Value {
	newMap := make(map[string]string)
	for key, value := range v {
		newMap[key] = value
	}

	return newMap
}

func copyItem(it *item) *item {
	return &item{
		key:   it.getKey(),
		value: copyValue(it.getValue()),
	}
}
