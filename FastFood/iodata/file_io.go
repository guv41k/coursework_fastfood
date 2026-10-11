package iodata

import (
	"bufio"
	"fmt"
	"os"

	"fastfood/managers"
	"fastfood/models/cashier"
	"fastfood/models/order"
)

func LoadCashiers(path string, cm *managers.CashierManager) (loaded, skipped int, details []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("file_io: %w", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := sc.Text()
		if line == "" {
			continue
		}
		c, err := cashier.FromString(line)
		if err != nil {
			skipped++
			details = append(details, fmt.Sprintf("строка %d: %v", lineNum, err))
			continue
		}
		if err := cm.Add(c); err != nil {
			skipped++
			details = append(details, fmt.Sprintf("строка %d: %v", lineNum, err))
			continue
		}
		loaded++
	}
	return loaded, skipped, details, nil
}

func SaveCashiers(path string, cm *managers.CashierManager) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("file_io: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < cm.Size(); i++ {
		c, _ := cm.At(i)
		fmt.Fprintln(w, c.String())
	}
	return w.Flush()
}

func LoadOrders(path string, om *managers.OrderManager, cm *managers.CashierManager) (loaded, skipped int, details []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("file_io: %w", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := sc.Text()
		if line == "" {
			continue
		}
		o, err := order.FromString(line)
		if err != nil {
			skipped++
			details = append(details, fmt.Sprintf("строка %d: %v", lineNum, err))
			continue
		}
		if err := om.Add(o, cm); err != nil {
			skipped++
			details = append(details, fmt.Sprintf("строка %d: %v", lineNum, err))
			continue
		}
		loaded++
	}
	return loaded, skipped, details, nil
}

func SaveOrders(path string, om *managers.OrderManager) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("file_io: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < om.Size(); i++ {
		o, _ := om.At(i)
		fmt.Fprintln(w, o.String())
	}
	return w.Flush()
}
