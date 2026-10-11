package managers

import (
	"fmt"

	"fastfood/ds/tree"
	"fastfood/models/date"
	"fastfood/models/order"
)

type OrderManager struct {
	items         []order.Order
	ByCashierAddr *tree.BST[string, int]
	byDate        *tree.BST[date.Date, int]
}

func dateLess(a, b date.Date) bool {
	if a.Year != b.Year {
		return a.Year < b.Year
	}
	if a.Month != b.Month {
		return a.Month < b.Month
	}
	return a.Day < b.Day
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		ByCashierAddr: tree.NewTree[string, int](func(a, b string) bool { return a < b }),
		byDate:        tree.NewTree[date.Date, int](dateLess),
	}
}

func orderKey(fio, address string) string {
	return fio + "|" + address
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
	m.ByCashierAddr.Insert(orderKey(o.CashierFIO, o.Address), idx)
	m.byDate.Insert(o.Date, idx)
	return nil
}

func (m *OrderManager) Remove(index int) bool {
	if index < 0 || index >= len(m.items) {
		return false
	}
	o := m.items[index]
	m.ByCashierAddr.Remove(orderKey(o.CashierFIO, o.Address), index)
	m.byDate.Remove(o.Date, index)

	last := len(m.items) - 1
	if index != last {
		moved := m.items[last]
		m.items[index] = moved
		m.ByCashierAddr.UpdateValue(orderKey(moved.CashierFIO, moved.Address), last, index)
		m.byDate.UpdateValue(moved.Date, last, index)
	}
	m.items = m.items[:last]
	return true
}

func (m *OrderManager) HasOrdersFor(fio, address string) bool {
	_, _, found := m.ByCashierAddr.Find(orderKey(fio, address))
	return found
}

func (m *OrderManager) FindByCashierAddress(fio, address string) ([]order.Order, int) {
	list, steps, found := m.ByCashierAddr.Find(orderKey(fio, address))
	if !found {
		return nil, steps
	}
	var result []order.Order
	for n := list.Head(); n != nil; n = n.Next {
		if o, ok := m.At(n.Value); ok {
			result = append(result, o)
		}
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
