package tree

import (
	"fmt"
	"strings"

	"fastfood/ds/linkedlist"
)

type node[K any, V comparable] struct {
	key   K
	list  *linkedlist.LinkedList[V]
	left  *node[K, V]
	right *node[K, V]
}

type BST[K any, V comparable] struct {
	root *node[K, V]
	less func(a, b K) bool
}

func NewTree[K any, V comparable](less func(a, b K) bool) *BST[K, V] {
	return &BST[K, V]{
		less: less,
	}
}

func (t *BST[K, V]) equal(a, b K) bool {
	return !t.less(a, b) && !t.less(b, a)
}

func (t *BST[K, V]) Insert(key K, value V) bool {
	var inserted bool
	t.root, inserted = t.insertNode(t.root, key, value)
	return inserted
}

func (t *BST[K, V]) insertNode(n *node[K, V], key K, value V) (*node[K, V], bool) {
	if n == nil {
		nn := &node[K, V]{
			key:  key,
			list: linkedlist.NewList[V](),
		}
		nn.list.Insert(value)
		return nn, true
	}
	inserted := true
	switch {
	case t.equal(key, n.key):
		for ln := n.list.Head(); ln != nil; ln = ln.Next {
			if ln.Value == value {
				return n, false
			}
		}
		n.list.Insert(value)
	case t.less(key, n.key):
		n.left, inserted = t.insertNode(n.left, key, value)
	default:
		n.right, inserted = t.insertNode(n.right, key, value)
	}
	return n, inserted
}

func (t *BST[K, V]) Find(key K) (list *linkedlist.LinkedList[V], steps int, found bool) {
	cur := t.root
	for cur != nil {
		steps++
		if t.equal(key, cur.key) {
			return cur.list, steps, true
		}
		if t.less(key, cur.key) {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return nil, steps, false
}

func (t *BST[K, V]) Remove(key K, value V) bool {
	var removed bool
	t.root, removed = t.removeFromTree(t.root, key, value)
	return removed
}

func (t *BST[K, V]) removeFromTree(n *node[K, V], key K, value V) (*node[K, V], bool) {
	if n == nil {
		return nil, false
	}
	var removed bool
	if t.less(key, n.key) {
		n.left, removed = t.removeFromTree(n.left, key, value)
		return n, removed
	}
	if t.less(n.key, key) {
		n.right, removed = t.removeFromTree(n.right, key, value)
		return n, removed
	}
	removed = n.list.Remove(func(v V) bool {
		return v == value
	})
	if n.list.Len() == 0 {
		return t.deleteNode(n), removed
	}
	return n, removed
}

func (t *BST[K, V]) deleteNode(n *node[K, V]) *node[K, V] {
	if n.left == nil {
		return n.right
	}
	if n.right == nil {
		return n.left
	}
	succ, succParent := t.findMin(n.right)
	n.key = succ.key
	n.list = succ.list
	if succParent == nil {
		n.right = succ.right
	} else {
		succParent.left = succ.right
	}
	return n
}

func (t *BST[K, V]) findMin(n *node[K, V]) (min *node[K, V], parent *node[K, V]) {
	for n.left != nil {
		parent = n
		n = n.left
	}
	return n, parent
}

func (t *BST[K, V]) UpdateValue(key K, oldValue, newValue V) bool {
	list, _, found := t.Find(key)
	if !found {
		return false
	}
	for n := list.Head(); n != nil; n = n.Next {
		if n.Value == oldValue {
			n.Value = newValue
			return true
		}
	}
	return false
}

func (t *BST[K, V]) PrintASCII(keyToString func(K) string) string {
	var sb strings.Builder
	t.printHelper(t.root, &sb, "", true, keyToString)
	return sb.String()
}

func (t *BST[K, V]) printHelper(n *node[K, V], sb *strings.Builder, prefix string, isLeft bool, keyToString func(K) string) {
	if n == nil {
		return
	}
	if n.right != nil {
		childPrefix := prefix
		if isLeft {
			childPrefix += "│   "
		} else {
			childPrefix += "    "
		}
		t.printHelper(n.right, sb, childPrefix, false, keyToString)
	}

	sb.WriteString(prefix)
	if isLeft {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("┌── ")
	}
	sb.WriteString(keyToString(n.key))
	fmt.Fprintf(sb, " [%d]\n", n.list.Len())

	if n.left != nil {
		childPrefix := prefix
		if isLeft {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
		t.printHelper(n.left, sb, childPrefix, true, keyToString)
	}
}

func (t *BST[K, V]) Clear() {
	t.root = nil
}
