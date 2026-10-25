package cashier

import (
	"fmt"
	"strconv"
	"strings"
)

type Cashier struct {
	FIO     string
	Address string
	Salary  float64
}

func isRuUpper(r rune) bool {
	return r == 'Ё' || (r >= 'А' && r <= 'Я')
}

func isRuLower(r rune) bool {
	return r == 'ё' || (r >= 'а' && r <= 'я')
}

func IsValidFIO(fio string) bool {
	words := strings.Split(fio, " ")
	if len(words) != 3 {
		return false
	}
	for _, w := range words {
		runes := []rune(w)
		if len(runes) == 0 || !isRuUpper(runes[0]) {
			return false
		}
		for _, r := range runes[1:] {
			if !isRuLower(r) {
				return false
			}
		}
	}
	return true
}

func (c Cashier) Validate() error {
	if !IsValidFIO(c.FIO) {
		return fmt.Errorf("cashier: invalid data(full name)")
	}
	if c.Address == "" || strings.Contains(c.Address, ";") {
		return fmt.Errorf("cashier: invalid data(address)")
	}
	if c.Salary <= 0 {
		return fmt.Errorf("cashier: invalid data(salary)")
	}
	return nil
}

func (c Cashier) String() string {
	return fmt.Sprintf("%s;%s;%.2f", c.FIO, c.Address, c.Salary)
}

func FromString(line string) (Cashier, error) {
	parts := strings.Split(line, ";")
	if len(parts) != 3 {
		return Cashier{}, fmt.Errorf("cashier: invalid data(field count)")
	}

	salary, _ := strconv.ParseFloat(parts[2], 64)
	return Cashier{FIO: parts[0], Address: parts[1], Salary: salary}, nil
}
