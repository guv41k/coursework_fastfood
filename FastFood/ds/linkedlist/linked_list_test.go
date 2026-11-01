package linkedlist

import (
	"slices"
	"testing"
)

func checkValues(t *testing.T, list *LinkedList[int], want []int) {
	t.Helper()
	var got []int
	for node := list.Head(); node != nil; node = node.Next {
		got = append(got, node.Value)
	}
	if list.Len() != len(want) || !slices.Equal(got, want) {
		t.Fatalf("list = %v (Len %d), want %v (Len %d)", got, list.Len(), want, len(want))
	}
}

func TestInsertAndRemovePositions(t *testing.T) {
	list := NewList[int]()
	checkValues(t, list, nil)
	if list.Remove(func(v int) bool { return v == 1 }) {
		t.Fatal("removed an item from an empty list")
	}

	for _, value := range []int{1, 2, 3} {
		list.Insert(value)
	}
	checkValues(t, list, []int{1, 2, 3})

	for _, step := range []struct {
		remove int
		want   []int
	}{
		{2, []int{1, 3}}, // middle
		{1, []int{3}},    // head
		{3, nil},         // tail and last node
	} {
		if !list.Remove(func(v int) bool { return v == step.remove }) {
			t.Fatalf("could not remove %d", step.remove)
		}
		checkValues(t, list, step.want)
	}

	list.Insert(4)
	checkValues(t, list, []int{4})
}

func TestRemoveOnlyFirstMatch(t *testing.T) {
	list := NewList[int]()
	for _, value := range []int{1, 2, 2, 3} {
		list.Insert(value)
	}
	if !list.Remove(func(v int) bool { return v == 2 }) {
		t.Fatal("could not remove matching item")
	}
	checkValues(t, list, []int{1, 2, 3})
	if list.Remove(func(v int) bool { return v == 9 }) {
		t.Fatal("removed a missing item")
	}
	checkValues(t, list, []int{1, 2, 3})
}
