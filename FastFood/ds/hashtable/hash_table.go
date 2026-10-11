package hashtable

import (
	"fmt"
	"strings"
)

type cellStatus int

const (
	cellEmpty cellStatus = iota
	cellOccupied
	cellDeleted
)

type cell[V any] struct {
	key    string
	value  V
	status cellStatus
}

type HashTable[V any] struct {
	cells []cell[V]
	size  int
	count int
}

func NewHashTable[V any](size int) *HashTable[V] {
	if size < 1 {
		size = 1
	}
	return &HashTable[V]{
		cells: make([]cell[V], size),
		size:  size,
	}
}

func hash(key string, size int) int {
	h := uint64(5381)
	for _, b := range []byte(key) {
		h = h*33 + uint64(b)
	}
	return int(h % uint64(size))
}

func (h *HashTable[V]) Insert(key string, value V) bool {
	if _, _, found := h.Find(key); found {
		return false
	}

	start := hash(key, h.size)
	firstDeleted := -1

	for i := 0; i < h.size; i++ {
		idx := (start + i) % h.size
		if h.cells[idx].status == cellEmpty {
			pos := idx
			if firstDeleted >= 0 {
				pos = firstDeleted
			}
			h.cells[pos] = cell[V]{
				key:    key,
				value:  value,
				status: cellOccupied,
			}
			h.count++
			return true
		}
		if h.cells[idx].status == cellDeleted && firstDeleted < 0 {
			firstDeleted = idx
		}
	}

	if firstDeleted >= 0 {
		h.cells[firstDeleted] = cell[V]{
			key:    key,
			value:  value,
			status: cellOccupied,
		}
		h.count++
		return true
	}
	return false
}

func (h *HashTable[V]) Find(key string) (value V, steps int, found bool) {
	start := hash(key, h.size)
	for i := 0; i < h.size; i++ {
		idx := (start + i) % h.size
		steps++
		switch h.cells[idx].status {
		case cellEmpty:
			return value, steps, false
		case cellDeleted:
			continue
		default:
			if h.cells[idx].key == key {
				return h.cells[idx].value, steps, true
			}
		}
	}
	return value, steps, false
}

func (h *HashTable[V]) Update(key string, newValue V) bool {
	start := hash(key, h.size)
	for i := 0; i < h.size; i++ {
		idx := (start + i) % h.size
		if h.cells[idx].status == cellEmpty {
			return false
		}
		if h.cells[idx].status == cellOccupied && h.cells[idx].key == key {
			h.cells[idx].value = newValue
			return true
		}
	}
	return false
}

func (h *HashTable[V]) Remove(key string) bool {
	start := hash(key, h.size)
	for i := 0; i < h.size; i++ {
		idx := (start + i) % h.size
		if h.cells[idx].status == cellEmpty {
			return false
		}
		if h.cells[idx].status == cellOccupied && h.cells[idx].key == key {
			h.cells[idx].status = cellDeleted
			h.count--
			return true
		}
	}
	return false
}

func (h *HashTable[V]) Print(valueToStr func(V) string) string {
	var sb strings.Builder
	for i := 0; i < h.size; i++ {
		fmt.Fprintf(&sb, "[%d] ", i)
		switch h.cells[i].status {
		case cellEmpty:
			sb.WriteString("статус=0 | —\n")
		case cellDeleted:
			sb.WriteString("статус=2 | DELETED\n")
		default:
			h1 := hash(h.cells[i].key, h.size)
			step := 0
			for p := 0; p < h.size; p++ {
				if (h1+p)%h.size == i {
					step = p
					break
				}
			}
			fmt.Fprintf(&sb, "статус=1 | h1=%d, шаг=%d, h2=%d | %s | %s\n",
				h1, step, i, h.cells[i].key, valueToStr(h.cells[i].value))
		}
	}
	fmt.Fprintf(&sb, "Заполнено: %d / %d\n", h.count, h.size)
	return sb.String()
}

func (h *HashTable[V]) Reset() {
	for i := range h.cells {
		h.cells[i] = cell[V]{}
	}
	h.count = 0
}
