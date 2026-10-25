package managers

import (
	"fmt"

	"fastfood/ds/tree"
	"fastfood/models/date"
	"fastfood/models/order"
)

type OrderManager struct {
	items         []order.Order
	byCashierAddr *tree.BST[string, int]
	byDate        *tree.BST[date.Date, int]
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		byCashierAddr: tree.NewTree[string, int](func(a, b string) bool { return a < b }),
		byDate:        tree.NewTree[date.Date, int](date.Date.Before),
	}
}

func (m *OrderManager) Add(o order.Order, cashiers *CashierManager) error {
	if err := o.Validate(); err != nil {
		return err
	}
	if !cashiers.Exists(o.CashierFIO, o.Address) {
		return fmt.Errorf("order: cashier not found")
	}
	m.items = append(m.items, o)
	idx := len(m.items) - 1
	m.byCashierAddr.Insert(fioAddrKey(o.CashierFIO, o.Address), idx)
	m.byDate.Insert(o.Date, idx)
	return nil
}

func (m *OrderManager) Remove(index int) bool {
	if index < 0 || index >= len(m.items) {
		return false
	}
	o := m.items[index]
	m.byCashierAddr.Remove(fioAddrKey(o.CashierFIO, o.Address), index)
	m.byDate.Remove(o.Date, index)

	last := len(m.items) - 1
	if index != last {
		moved := m.items[last]
		m.items[index] = moved
		m.byCashierAddr.UpdateValue(fioAddrKey(moved.CashierFIO, moved.Address), last, index)
		m.byDate.UpdateValue(moved.Date, last, index)
	}
	m.items = m.items[:last]
	return true
}

func (m *OrderManager) HasOrdersFor(fio, address string) bool {
	_, _, found := m.byCashierAddr.Find(fioAddrKey(fio, address))
	return found
}

func (m *OrderManager) FindByCashierAddress(fio, address string) ([]order.Order, int) {
	list, steps, found := m.byCashierAddr.Find(fioAddrKey(fio, address))
	if !found {
		return nil, steps
	}
	var result []order.Order
	for n := list.Head(); n != nil; n = n.Next {
		result = append(result, m.items[n.Value])
	}
	return result, steps
}

func (m *OrderManager) Size() int {
	return len(m.items)
}

func (m *OrderManager) At(i int) (order.Order, bool) {
	if i < 0 || i >= len(m.items) {
		return order.Order{}, false
	}
	return m.items[i], true
}

func (m *OrderManager) Reset() {
	m.items = nil
	m.byCashierAddr.Clear()
	m.byDate.Clear()
}

func (m *OrderManager) PrintByCashierAddr() string {
	return m.byCashierAddr.PrintASCII(func(k string) string { return k })
}

func (m *OrderManager) PrintByDate() string {
	return m.byDate.PrintASCII(date.Date.String)
}
