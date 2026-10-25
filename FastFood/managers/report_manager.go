package managers

import (
	"math"

	"fastfood/models/date"
)

type ReportRow struct {
	CashierFIO     string
	CashierAddress string
	CashierSalary  float64
	OrderDate      string
	OrderTime      string
	OrderSum       float64
}

func BuildReport(cashiers *CashierManager, orders *OrderManager, address string, d date.Date, sum float64) []ReportRow {
	list, _, found := orders.byDate.Find(d)
	if !found {
		return nil
	}

	var rows []ReportRow
	for n := list.Head(); n != nil; n = n.Next {
		o := orders.items[n.Value]
		if o.Address != address || math.Abs(o.Sum-sum) >= 0.005 {
			continue
		}
		c, _, _ := cashiers.Find(o.CashierFIO, o.Address)
		rows = append(rows, ReportRow{
			CashierFIO:     c.FIO,
			CashierAddress: c.Address,
			CashierSalary:  c.Salary,
			OrderDate:      o.Date.String(),
			OrderTime:      o.Time,
			OrderSum:       o.Sum,
		})
	}
	return rows
}
