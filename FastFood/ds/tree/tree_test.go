package tree

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func newIntTree() *BST[int, int] {
	return NewTree[int, int](func(a, b int) bool { return a < b })
}

func checkTreeValues(t *testing.T, tree *BST[int, int], key int, want []int) {
	t.Helper()
	list, _, found := tree.Find(key)
	if !found {
		t.Fatalf("key %d not found", key)
	}
	var got []int
	for node := list.Head(); node != nil; node = node.Next {
		got = append(got, node.Value)
	}
	if !slices.Equal(got, want) || list.Len() != len(want) {
		t.Fatalf("values for key %d = %v (Len %d), want %v", key, got, list.Len(), want)
	}
}

func TestInsertFindUpdateAndRemoveValues(t *testing.T) {
	tree := newIntTree()
	if _, steps, found := tree.Find(10); found || steps != 0 {
		t.Fatalf("Find on empty tree = (steps %d, found %t), want (0, false)", steps, found)
	}
	for _, entry := range []struct{ key, value int }{{10, 1}, {5, 2}, {15, 3}, {10, 4}} {
		if !tree.Insert(entry.key, entry.value) {
			t.Fatalf("could not insert (%d, %d)", entry.key, entry.value)
		}
	}
	if tree.Insert(10, 1) {
		t.Fatal("duplicate key/value pair was inserted")
	}
	checkTreeValues(t, tree, 10, []int{1, 4})
	if _, steps, found := tree.Find(5); !found || steps != 2 {
		t.Fatalf("Find(5) = (steps %d, found %t), want (2, true)", steps, found)
	}
	if _, _, found := tree.Find(99); found {
		t.Fatal("found a missing key")
	}

	if !tree.UpdateValue(10, 4, 9) || tree.UpdateValue(10, 4, 8) || tree.UpdateValue(99, 1, 2) {
		t.Fatal("UpdateValue returned an unexpected result")
	}
	checkTreeValues(t, tree, 10, []int{1, 9})
	if tree.Remove(10, 4) || !tree.Remove(10, 1) {
		t.Fatal("Remove returned an unexpected result")
	}
	if tree.Remove(99, 1) {
		t.Fatal("removed a value from a missing key")
	}
	checkTreeValues(t, tree, 10, []int{9})
	if !tree.Remove(10, 9) {
		t.Fatal("could not remove the last value for a key")
	}
	if _, _, found := tree.Find(10); found {
		t.Fatal("key remained after its last value was removed")
	}
	checkTreeValues(t, tree, 5, []int{2})
	checkTreeValues(t, tree, 15, []int{3})
}

func TestRemoveNodeShapes(t *testing.T) {
	tests := []struct {
		name     string
		keys     []int
		remove   int
		wantRoot int
	}{
		{"leaf", []int{5, 3, 7}, 3, 5},
		{"one child", []int{5, 3, 2}, 3, 5},
		{"two children, immediate successor", []int{5, 3, 7, 8}, 5, 7},
		{"two children, deep successor with right child", []int{10, 5, 20, 15, 25, 12, 17, 13}, 10, 12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := newIntTree()
			for _, key := range tt.keys {
				if !tree.Insert(key, key*10) {
					t.Fatalf("could not insert key %d", key)
				}
			}
			if !tree.Remove(tt.remove, tt.remove*10) {
				t.Fatalf("could not remove key %d", tt.remove)
			}
			if tree.root == nil || tree.root.key != tt.wantRoot {
				t.Fatalf("root after removal = %v, want key %d", tree.root, tt.wantRoot)
			}
			if _, _, found := tree.Find(tt.remove); found {
				t.Fatalf("removed key %d is still present", tt.remove)
			}
			for _, key := range tt.keys {
				if key != tt.remove {
					checkTreeValues(t, tree, key, []int{key * 10})
				}
			}
		})
	}
}

func TestClear(t *testing.T) {
	tree := newIntTree()
	tree.Insert(1, 10)
	tree.Clear()
	if _, _, found := tree.Find(1); found || tree.root != nil {
		t.Fatal("Clear left the old node accessible")
	}
	if !tree.Insert(2, 20) {
		t.Fatal("tree cannot be reused after Clear")
	}
	checkTreeValues(t, tree, 2, []int{20})
}

func TestRemoveLastNode(t *testing.T) {
	tree := newIntTree()
	tree.Insert(1, 10)
	if !tree.Remove(1, 10) || tree.root != nil {
		t.Fatal("removing the last value did not empty the tree")
	}
}

func TestPrintASCIIShowsKeysAndValueCounts(t *testing.T) {
	tree := newIntTree()
	for _, entry := range []struct{ key, value int }{{5, 1}, {3, 2}, {7, 3}, {5, 4}} {
		tree.Insert(entry.key, entry.value)
	}

	output := tree.PrintASCII(strconv.Itoa)
	for _, part := range []string{"5 [2]", "3 [1]", "7 [1]"} {
		if !strings.Contains(output, part) {
			t.Errorf("PrintASCII output %q does not contain %q", output, part)
		}
	}
}
