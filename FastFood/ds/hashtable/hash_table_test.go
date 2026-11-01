package hashtable

import (
	"strconv"
	"strings"
	"testing"
)

func TestCollisionDeletionAndReuse(t *testing.T) {
	const size = 3
	if hash("a", size) != hash("d", size) || hash("a", size) != hash("g", size) {
		t.Fatal("test keys must collide")
	}

	table := NewHashTable[int](size)
	if !table.Insert("a", 10) || !table.Insert("d", 20) {
		t.Fatal("could not insert colliding keys")
	}
	if value, steps, found := table.Find("d"); !found || value != 20 || steps != 2 {
		t.Fatalf("Find(d) = (%d, %d, %t), want (20, 2, true)", value, steps, found)
	}
	if table.Insert("d", 99) {
		t.Fatal("duplicate key was inserted")
	}

	if !table.Remove("a") {
		t.Fatal("could not remove first key in probe chain")
	}
	if table.Insert("d", 99) {
		t.Fatal("duplicate key was inserted after a tombstone")
	}
	if value, steps, found := table.Find("d"); !found || value != 20 || steps != 2 {
		t.Fatalf("Find(d) after deletion = (%d, %d, %t), want (20, 2, true)", value, steps, found)
	}
	if !table.Insert("g", 30) {
		t.Fatal("could not reuse deleted cell")
	}
	if value, steps, found := table.Find("g"); !found || value != 30 || steps != 1 {
		t.Fatalf("Find(g) = (%d, %d, %t), want (30, 1, true)", value, steps, found)
	}

	if !table.Update("d", 22) || table.Update("missing", 1) {
		t.Fatal("Update returned an unexpected result")
	}
	if value, _, found := table.Find("d"); !found || value != 22 {
		t.Fatalf("Find(d) after update = (%d, %t), want (22, true)", value, found)
	}
	if !table.Insert("j", 40) || table.Insert("m", 50) {
		t.Fatal("full table accepted an extra key or rejected the last free cell")
	}
	if _, steps, found := table.Find("m"); found || steps != size {
		t.Fatalf("Find(m) in full table = (steps %d, found %t), want (%d, false)", steps, found, size)
	}
}

func TestResetAndMinimumSize(t *testing.T) {
	table := NewHashTable[int](0)
	if !table.Insert("first", 1) || table.Insert("second", 2) {
		t.Fatal("minimum-size table should hold exactly one key")
	}

	table.Reset()
	if _, _, found := table.Find("first"); found {
		t.Fatal("Reset left the old key accessible")
	}
	if !table.Insert("second", 2) {
		t.Fatal("table cannot be reused after Reset")
	}
	if table.Remove("first") {
		t.Fatal("Remove accepted a missing key")
	}
}

func TestPrintShowsCellStates(t *testing.T) {
	table := NewHashTable[int](3)
	table.Insert("a", 10)
	table.Insert("d", 20)
	table.Remove("a")

	output := table.Print(strconv.Itoa)
	for _, part := range []string{"статус=0", "DELETED", "d | 20", "Заполнено: 1 / 3"} {
		if !strings.Contains(output, part) {
			t.Errorf("Print output %q does not contain %q", output, part)
		}
	}
}
