package managers

import (
	"fmt"

	"fastfood/ds/hashtable"
	"fastfood/models/cashier"
)

type CashierManager struct {
	items []cashier.Cashier
	index *hashtable.HashTable[int]
}

func fioAddrKey(fio, address string) string {
	return fio + "|" + address
}

func NewCashierManager(hashSize int) *CashierManager {
	return &CashierManager{
		index: hashtable.NewHashTable[int](hashSize),
	}
}

func (m *CashierManager) Add(c cashier.Cashier) error {
	if err := c.Validate(); err != nil {
		return err
	}
	key := fioAddrKey(c.FIO, c.Address)
	if _, _, found := m.index.Find(key); found {
		return fmt.Errorf("cashier: already exists")
	}
	m.items = append(m.items, c)
	if !m.index.Insert(key, len(m.items)-1) {
		m.items = m.items[:len(m.items)-1]
		return fmt.Errorf("cashier: hash table full")
	}
	return nil
}

func (m *CashierManager) Remove(fio, address string, orders *OrderManager) error {
	if orders.HasOrdersFor(fio, address) {
		return fmt.Errorf("cashier: has orders")
	}
	key := fioAddrKey(fio, address)
	idx, _, found := m.index.Find(key)
	if !found {
		return fmt.Errorf("cashier: not found")
	}
	m.index.Remove(key)

	last := len(m.items) - 1
	if idx != last {
		moved := m.items[last]
		m.items[idx] = moved
		m.index.Update(fioAddrKey(moved.FIO, moved.Address), idx)
	}
	m.items = m.items[:last]
	return nil
}

func (m *CashierManager) Find(fio, address string) (cashier.Cashier, int, bool) {
	idx, steps, found := m.index.Find(fioAddrKey(fio, address))
	if !found {
		return cashier.Cashier{}, steps, false
	}
	return m.items[idx], steps, true
}

func (m *CashierManager) Exists(fio, address string) bool {
	_, _, found := m.Find(fio, address)
	return found
}

func (m *CashierManager) Size() int {
	return len(m.items)
}

func (m *CashierManager) At(i int) (cashier.Cashier, bool) {
	if i < 0 || i >= len(m.items) {
		return cashier.Cashier{}, false
	}
	return m.items[i], true
}

func (m *CashierManager) Reset() {
	m.items = nil
	m.index.Reset()
}

func (m *CashierManager) Print() string {
	return m.index.Print(func(i int) string { return m.items[i].String() })
}
