package iodata

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"fastfood/managers"
	"fastfood/models/cashier"
	"fastfood/models/order"
)

func LoadCashiers(r io.Reader, cm *managers.CashierManager) (loaded, skipped int, details []string, err error) {
	sc := bufio.NewScanner(r)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "" {
			continue
		}
		c, err := cashier.FromString(line)
		if err == nil {
			err = cm.Add(c)
		}
		if err != nil {
			skipped++
			details = append(details, fmt.Sprintf("строка %d: %v", lineNum, err))
			continue
		}
		loaded++
	}
	if err := sc.Err(); err != nil {
		return loaded, skipped, details, fmt.Errorf("file_io: %w", err)
	}
	return loaded, skipped, details, nil
}

func SaveCashiers(w io.Writer, cm *managers.CashierManager) error {
	bw := bufio.NewWriter(w)
	for i := 0; i < cm.Size(); i++ {
		c, _ := cm.At(i)
		if _, err := fmt.Fprintln(bw, c.String()); err != nil {
			return err
		}
	}
	return bw.Flush()
}

func LoadOrders(r io.Reader, om *managers.OrderManager, cm *managers.CashierManager) (loaded, skipped int, details []string, err error) {
	sc := bufio.NewScanner(r)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "" {
			continue
		}
		o, err := order.FromString(line)
		if err == nil {
			err = om.Add(o, cm)
		}
		if err != nil {
			skipped++
			details = append(details, fmt.Sprintf("строка %d: %v", lineNum, err))
			continue
		}
		loaded++
	}
	if err := sc.Err(); err != nil {
		return loaded, skipped, details, fmt.Errorf("file_io: %w", err)
	}
	return loaded, skipped, details, nil
}

func SaveOrders(w io.Writer, om *managers.OrderManager) error {
	bw := bufio.NewWriter(w)
	for i := 0; i < om.Size(); i++ {
		o, _ := om.At(i)
		if _, err := fmt.Fprintln(bw, o.String()); err != nil {
			return err
		}
	}
	return bw.Flush()
}
