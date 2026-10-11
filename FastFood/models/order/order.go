package order

import (
	"fmt"
	"strconv"
	"strings"

	"fastfood/models/cashier"
	"fastfood/models/date"
)

type Order struct {
	CashierFIO string
	Address    string
	Date       date.Date
	Time       string
	Sum        float64
}

func isValidTime(t string) bool {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return false
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	return err1 == nil && err2 == nil && h >= 0 && h <= 23 && m >= 0 && m <= 59
}

func (o Order) Validate() error {
	if !cashier.IsValidFIO(o.CashierFIO) {
		return fmt.Errorf("order: invalid data(cashier full name)")
	}
	if o.Address == "" || strings.Contains(o.Address, ";") {
		return fmt.Errorf("order: invalid data(address)")
	}
	if !o.Date.IsValid() {
		return fmt.Errorf("order: invalid data(date)")
	}
	if !isValidTime(o.Time) {
		return fmt.Errorf("order: invalid data(time)")
	}
	if o.Sum <= 0 {
		return fmt.Errorf("order: invalid data(sum)")
	}
	return nil
}

func (o Order) String() string {
	return fmt.Sprintf("%s;%s;%s;%s;%.2f", o.CashierFIO, o.Address, o.Date.String(), o.Time, o.Sum)
}

func FromString(line string) (Order, error) {
	parts := strings.Split(line, ";")
	if len(parts) != 5 {
		return Order{}, fmt.Errorf("order: invalid data(field count)")
	}

	d, _ := date.Parse(parts[2])
	sum, _ := strconv.ParseFloat(parts[4], 64)
	o := Order{
		CashierFIO: parts[0],
		Address:    parts[1],
		Date:       d,
		Time:       parts[3],
		Sum:        sum,
	}
	if err := o.Validate(); err != nil {
		return Order{}, err
	}
	return o, nil
}
