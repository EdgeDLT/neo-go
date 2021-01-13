package mpt

import (
	"bytes"
	"sort"
)

// Batch is batch of storage changes.
// It stores key-value pairs in a sorted state.
type Batch struct {
	kv []keyValue
}

type keyValue struct {
	key   []byte
	value []byte
}

// Add adds key-value pair to batch.
// If there is an item with the specified key, it is replaced.
func (b *Batch) Add(key []byte, value []byte) {
	path := toNibbles(key)
	i := sort.Search(len(b.kv), func(i int) bool {
		return bytes.Compare(path, b.kv[i].key) <= 0
	})
	if i == len(b.kv) {
		b.kv = append(b.kv, keyValue{path, value})
	} else if bytes.Equal(b.kv[i].key, path) {
		b.kv[i].value = value
	} else {
		b.kv = append(b.kv, keyValue{})
		copy(b.kv[i+1:], b.kv[i:])
		b.kv[i].key = path
		b.kv[i].value = value
	}
}

// PutBatch puts batch to trie.
// It is not atomic (and probably cannot be without substantial slow-down)
// and returns number of elements processed.
// However each element is being put atomically, so Trie is always in a valid state.
// It is used mostly after the block processing to update MPT and error is not expected.
func (t *Trie) PutBatch(b Batch) (int, error) {
	r, n, err := t.putBatch(b.kv)
	t.root = r
	return n, err
}

func (t *Trie) putBatch(kv []keyValue) (Node, int, error) {
	return t.putBatchIntoNode(t.root, kv)
}

func (t *Trie) putBatchIntoNode(curr Node, kv []keyValue) (Node, int, error) {
	switch n := curr.(type) {
	case *LeafNode:
		return t.putBatchIntoLeaf(n, kv)
	case *BranchNode:
		return t.putBatchIntoBranch(n, kv)
	case *ExtensionNode:
		return t.putBatchIntoExtension(n, kv)
	case *HashNode:
		return t.putBatchIntoHash(n, kv)
	default:
		panic("invalid MPT node type")
	}
}

func (t *Trie) putBatchIntoLeaf(curr *LeafNode, kv []keyValue) (Node, int, error) {
	t.removeRef(curr.Hash(), curr.Bytes())
	return t.newSubTrieMany(nil, kv, curr.value)
}

func (t *Trie) putBatchIntoBranch(curr *BranchNode, kv []keyValue) (Node, int, error) {
	return t.addToBranch(curr, kv, true)
}

func (t *Trie) mergeExtension(prefix []byte, sub Node) Node {
	switch sn := sub.(type) {
	case *ExtensionNode:
		t.removeRef(sn.Hash(), sn.bytes)
		sn.key = append(prefix, sn.key...)
		sn.invalidateCache()
		t.addRef(sn.Hash(), sn.bytes)
		return sn
	case *HashNode:
		return sn
	default:
		if len(prefix) != 0 {
			e := NewExtensionNode(prefix, sub)
			t.addRef(e.Hash(), e.bytes)
			return e
		}
		return sub
	}
}

func (t *Trie) putBatchIntoExtension(curr *ExtensionNode, kv []keyValue) (Node, int, error) {
	t.removeRef(curr.Hash(), curr.bytes)

	common := lcpMany(kv)
	pref := lcp(common, curr.key)
	if len(pref) == len(curr.key) {
		// Extension must be split into new nodes.
		stripPrefix(len(curr.key), kv)
		sub, n, err := t.putBatchIntoNode(curr.next, kv)
		return t.mergeExtension(pref, sub), n, err
	}

	if len(pref) != 0 {
		stripPrefix(len(pref), kv)
		sub, n, err := t.putBatchIntoExtensionNoPrefix(curr.key[len(pref):], curr.next, kv)
		return t.mergeExtension(pref, sub), n, err
	}
	return t.putBatchIntoExtensionNoPrefix(curr.key, curr.next, kv)
}

func (t *Trie) putBatchIntoExtensionNoPrefix(key []byte, next Node, kv []keyValue) (Node, int, error) {
	b := NewBranchNode()
	if len(key) > 1 {
		b.Children[key[0]] = t.newSubTrie(key[1:], next, false)
	} else {
		b.Children[key[0]] = next
	}
	return t.addToBranch(b, kv, false)
}

func isEmpty(n Node) bool {
	hn, ok := n.(*HashNode)
	return ok && hn.IsEmpty()
}

// addToBranch puts items into the branch node assuming b is not yet in trie.
func (t *Trie) addToBranch(b *BranchNode, kv []keyValue, inTrie bool) (Node, int, error) {
	if inTrie {
		t.removeRef(b.Hash(), b.bytes)
	}
	n, err := t.iterateBatch(kv, func(c byte, kv []keyValue) (int, error) {
		child, n, err := t.putBatchIntoNode(b.Children[c], kv)
		b.Children[c] = child
		return n, err
	})
	if inTrie && n != 0 {
		b.invalidateCache()
	}
	return t.stripBranch(b), n, err
}

// stripsBranch strips branch node after incomplete batch put.
// It assumes there is no reference to b in trie.
func (t *Trie) stripBranch(b *BranchNode) Node {
	var n int
	var lastIndex byte
	for i := range b.Children {
		if !isEmpty(b.Children[i]) {
			n++
			lastIndex = byte(i)
		}
	}
	switch {
	case n == 0:
		return new(HashNode)
	case n == 1:
		return t.mergeExtension([]byte{lastIndex}, b.Children[lastIndex])
	default:
		t.addRef(b.Hash(), b.bytes)
		return b
	}
}

func (t *Trie) iterateBatch(kv []keyValue, f func(c byte, kv []keyValue) (int, error)) (int, error) {
	var n int
	for len(kv) != 0 {
		c, i := getLastIndex(kv)
		if c != lastChild {
			stripPrefix(1, kv[:i])
		}
		sub, err := f(c, kv[:i])
		n += sub
		if err != nil {
			return n, err
		}
		kv = kv[i:]
	}
	return n, nil
}

func (t *Trie) putBatchIntoHash(curr *HashNode, kv []keyValue) (Node, int, error) {
	if curr.IsEmpty() {
		common := lcpMany(kv)
		stripPrefix(len(common), kv)
		return t.newSubTrieMany(common, kv, nil)
	}
	result, err := t.getFromStore(curr.hash)
	if err != nil {
		return curr, 0, err
	}
	return t.putBatchIntoNode(result, kv)
}

// Creates new subtrie from provided key-value pairs.
// Items in kv must have no common prefix.
// If there are any deletions in kv, return error.
// kv is not empty.
// kv is sorted by key.
// value is current value stored by prefix.
func (t *Trie) newSubTrieMany(prefix []byte, kv []keyValue, value []byte) (Node, int, error) {
	if len(kv[0].key) == 0 {
		if len(kv[0].value) == 0 {
			if len(kv) == 1 {
				if len(value) != 0 {
					return new(HashNode), 1, nil
				}
				return new(HashNode), 0, ErrNotFound
			}
			node, n, err := t.newSubTrieMany(prefix, kv[1:], nil)
			return node, n + 1, err
		}
		if len(kv) == 1 {
			return t.newSubTrie(prefix, NewLeafNode(kv[0].value), true), 1, nil
		}
		value = kv[0].value
	}

	// Prefix is empty and we have at least 2 children.
	b := NewBranchNode()
	if len(value) != 0 {
		// Empty key is always first.
		leaf := NewLeafNode(value)
		t.addRef(leaf.Hash(), leaf.bytes)
		b.Children[lastChild] = leaf
	}
	nd, n, err := t.addToBranch(b, kv, false)
	return t.mergeExtension(prefix, nd), n, err
}

func stripPrefix(n int, kv []keyValue) {
	for i := range kv {
		kv[i].key = kv[i].key[n:]
	}
}

func getLastIndex(kv []keyValue) (byte, int) {
	if len(kv[0].key) == 0 {
		return lastChild, 1
	}
	c := kv[0].key[0]
	for i := range kv[1:] {
		if kv[i+1].key[0] != c {
			return c, i + 1
		}
	}
	return c, len(kv)
}
