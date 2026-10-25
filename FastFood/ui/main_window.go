package ui

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fastfood/iodata"
	"fastfood/managers"
	"fastfood/models/cashier"
	"fastfood/models/date"
	"fastfood/models/order"
)

type MainWindow struct {
	win fyne.Window
	cm  *managers.CashierManager
	om  *managers.OrderManager

	cashierRows *fyne.Container
	orderRows   *fyne.Container

	selectedCashierRow int
	selectedOrderRow   int
}

func NewMainWindow(win fyne.Window, hashSize int) *MainWindow {
	mw := &MainWindow{
		win:                win,
		cm:                 managers.NewCashierManager(hashSize),
		om:                 managers.NewOrderManager(),
		selectedCashierRow: -1,
		selectedOrderRow:   -1,
	}
	mw.build()
	return mw
}

func (mw *MainWindow) refresh() {
	mw.refreshCashierTable()
	mw.refreshOrderTable()
}

func (mw *MainWindow) uniqueFioAddr() ([]string, []string) {
	var fios, addrs []string
	seenFio := map[string]bool{}
	seenAddr := map[string]bool{}
	for i := 0; i < mw.cm.Size(); i++ {
		c, _ := mw.cm.At(i)
		if !seenFio[c.FIO] {
			seenFio[c.FIO] = true
			fios = append(fios, c.FIO)
		}
		if !seenAddr[c.Address] {
			seenAddr[c.Address] = true
			addrs = append(addrs, c.Address)
		}
	}
	return fios, addrs
}

func (mw *MainWindow) build() {
	mw.cashierRows = container.NewVBox()
	mw.orderRows = container.NewVBox()
	mw.refresh()

	cashierHeader := widget.NewLabelWithStyle("Кассиры (ПКМ на строке — меню операций)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	cashierPanel := container.NewBorder(cashierHeader, nil, nil, nil, container.NewVScroll(mw.cashierRows))

	orderHeader := widget.NewLabelWithStyle("Заказы (ПКМ на строке — меню операций)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	orderPanel := container.NewBorder(orderHeader, nil, nil, nil, container.NewVScroll(mw.orderRows))

	split := container.NewVSplit(cashierPanel, orderPanel)
	split.Offset = 0.45

	mw.win.SetContent(split)
	mw.win.Resize(fyne.NewSize(1050, 700))
}

func (mw *MainWindow) showCashierMenu(pos fyne.Position) {
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Добавить", mw.onAddCashier),
		fyne.NewMenuItem("Удалить", mw.onDeleteCashier),
		fyne.NewMenuItem("Найти", mw.onFindCashier),
		fyne.NewMenuItem("Показать структуры", mw.onShowStructures),
		fyne.NewMenuItem("Загрузить...", mw.onLoadCashiers),
		fyne.NewMenuItem("Сохранить...", mw.onSaveCashiers),
	)
	widget.NewPopUpMenu(menu, mw.win.Canvas()).ShowAtPosition(pos)
}

func (mw *MainWindow) showOrderMenu(pos fyne.Position) {
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Добавить", mw.onAddOrder),
		fyne.NewMenuItem("Удалить", mw.onDeleteOrder),
		fyne.NewMenuItem("Найти", mw.onFindOrders),
		fyne.NewMenuItem("Отчёт", mw.onBuildReport),
		fyne.NewMenuItem("Показать структуры", mw.onShowStructures),
		fyne.NewMenuItem("Загрузить...", mw.onLoadOrders),
		fyne.NewMenuItem("Сохранить...", mw.onSaveOrders),
	)
	widget.NewPopUpMenu(menu, mw.win.Canvas()).ShowAtPosition(pos)
}

func (mw *MainWindow) refreshCashierTable() {
	mw.cashierRows.RemoveAll()

	header := newTableRow(3, nil, func(_ int, pos fyne.Position) { mw.showCashierMenu(pos) })
	header.SetHeader("ФИО", "Адрес точки", "Оклад, руб.")
	mw.cashierRows.Add(header)

	for i := 0; i < mw.cm.Size(); i++ {
		c, _ := mw.cm.At(i)
		row := newTableRow(3,
			func(r int) { mw.selectedCashierRow = r; mw.refreshCashierTable() },
			func(r int, pos fyne.Position) { mw.selectedCashierRow = r; mw.showCashierMenu(pos) },
		)
		row.row = i
		row.SetTexts(i == mw.selectedCashierRow, c.FIO, c.Address, fmt.Sprintf("%.2f", c.Salary))
		mw.cashierRows.Add(row)
	}
	mw.cashierRows.Refresh()
}

func (mw *MainWindow) refreshOrderTable() {
	mw.orderRows.RemoveAll()

	header := newTableRow(5, nil, func(_ int, pos fyne.Position) { mw.showOrderMenu(pos) })
	header.SetHeader("ФИО кассира", "Адрес точки", "Дата", "Время", "Сумма, руб.")
	mw.orderRows.Add(header)

	for i := 0; i < mw.om.Size(); i++ {
		o, _ := mw.om.At(i)
		row := newTableRow(5,
			func(r int) { mw.selectedOrderRow = r; mw.refreshOrderTable() },
			func(r int, pos fyne.Position) { mw.selectedOrderRow = r; mw.showOrderMenu(pos) },
		)
		row.row = i
		row.SetTexts(i == mw.selectedOrderRow, o.CashierFIO, o.Address, o.Date.String(), o.Time, fmt.Sprintf("%.2f", o.Sum))
		mw.orderRows.Add(row)
	}
	mw.orderRows.Refresh()
}

func (mw *MainWindow) onAddCashier() {
	fioEntry := widget.NewEntry()
	addrEntry := widget.NewEntry()
	salaryEntry := widget.NewEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("ФИО кассира", fioEntry),
		widget.NewFormItem("Адрес точки", addrEntry),
		widget.NewFormItem("Оклад, руб.", salaryEntry),
	}
	dialog.ShowForm("Добавить кассира", "Добавить", "Отмена", items, func(ok bool) {
		if !ok {
			return
		}
		salary, _ := strconv.ParseFloat(strings.TrimSpace(salaryEntry.Text), 64)
		c := cashier.Cashier{
			FIO:     strings.TrimSpace(fioEntry.Text),
			Address: strings.TrimSpace(addrEntry.Text),
			Salary:  salary,
		}
		if err := mw.cm.Add(c); err != nil {
			dialog.ShowError(err, mw.win)
			return
		}
		mw.refresh()
	}, mw.win)
}

func (mw *MainWindow) onDeleteCashier() {
	if mw.selectedCashierRow < 0 {
		dialog.ShowInformation("Удаление кассира", "Строка не выбрана", mw.win)
		return
	}
	c, ok := mw.cm.At(mw.selectedCashierRow)
	if !ok {
		return
	}
	dialog.ShowConfirm("Подтверждение", fmt.Sprintf("Удалить кассира «%s»?", c.FIO), func(ok bool) {
		if !ok {
			return
		}
		if err := mw.cm.Remove(c.FIO, c.Address, mw.om); err != nil {
			dialog.ShowError(err, mw.win)
			return
		}
		mw.selectedCashierRow = -1
		mw.refresh()
	}, mw.win)
}

func (mw *MainWindow) onFindCashier() {
	if mw.cm.Size() == 0 {
		dialog.ShowInformation("Поиск кассира", "Нет кассиров для поиска", mw.win)
		return
	}
	fios, addrs := mw.uniqueFioAddr()
	fioSelect := widget.NewSelect(fios, nil)
	addrSelect := widget.NewSelect(addrs, nil)
	items := []*widget.FormItem{
		widget.NewFormItem("ФИО кассира", fioSelect),
		widget.NewFormItem("Адрес точки", addrSelect),
	}
	dialog.ShowForm("Найти кассира", "Найти", "Отмена", items, func(ok bool) {
		if !ok {
			return
		}
		c, steps, found := mw.cm.Find(fioSelect.Selected, addrSelect.Selected)
		if !found {
			dialog.ShowInformation("Поиск кассира",
				fmt.Sprintf("Кассир не найден\nШагов поиска в ХТ: %d", steps), mw.win)
			return
		}
		dialog.ShowInformation("Кассир найден",
			fmt.Sprintf("%s\n%s\n%.2f руб.\n\nШагов поиска в ХТ: %d", c.FIO, c.Address, c.Salary, steps), mw.win)
	}, mw.win)
}

func (mw *MainWindow) onShowStructures() {
	monoTab := func(title, text string) *container.TabItem {
		label := widget.NewLabel(text)
		label.TextStyle = fyne.TextStyle{Monospace: true}
		return container.NewTabItem(title, container.NewScroll(label))
	}

	content := container.NewAppTabs(
		monoTab("Хеш-таблица кассиров", mw.cm.Print()),
		monoTab("Дерево заказов (ФИО|адрес)", mw.om.PrintByCashierAddr()),
		monoTab("Дерево заказов (дата)", mw.om.PrintByDate()),
	)
	d := dialog.NewCustom("Структуры данных", "Закрыть", content, mw.win)
	d.Resize(fyne.NewSize(600, 600))
	d.Show()
}

func (mw *MainWindow) onAddOrder() {
	if mw.cm.Size() == 0 {
		dialog.ShowInformation("Добавить заказ", "Нет ни одного кассира", mw.win)
		return
	}
	var items []string
	for i := 0; i < mw.cm.Size(); i++ {
		c, _ := mw.cm.At(i)
		items = append(items, c.FIO+" | "+c.Address)
	}
	cashierSelect := widget.NewSelect(items, nil)
	dateEntry := widget.NewEntry()
	dateEntry.SetText("01.01.2025")
	timeEntry := widget.NewEntry()
	timeEntry.SetText("10:00")
	sumEntry := widget.NewEntry()
	sumEntry.SetText("300.00")

	formItems := []*widget.FormItem{
		widget.NewFormItem("Кассир", cashierSelect),
		widget.NewFormItem("Дата (ДД.ММ.ГГГГ)", dateEntry),
		widget.NewFormItem("Время (ЧЧ:ММ)", timeEntry),
		widget.NewFormItem("Сумма, руб.", sumEntry),
	}
	dialog.ShowForm("Добавить заказ", "Добавить", "Отмена", formItems, func(ok bool) {
		if !ok {
			return
		}
		idx := -1
		for i, s := range items {
			if s == cashierSelect.Selected {
				idx = i
				break
			}
		}
		if idx < 0 {
			return
		}
		c, _ := mw.cm.At(idx)

		d, _ := date.Parse(strings.TrimSpace(dateEntry.Text))
		sum, _ := strconv.ParseFloat(strings.TrimSpace(sumEntry.Text), 64)
		o := order.Order{
			CashierFIO: c.FIO,
			Address:    c.Address,
			Date:       d,
			Time:       strings.TrimSpace(timeEntry.Text),
			Sum:        sum,
		}
		if err := mw.om.Add(o, mw.cm); err != nil {
			dialog.ShowError(err, mw.win)
			return
		}
		mw.refresh()
	}, mw.win)
}

func (mw *MainWindow) onDeleteOrder() {
	if mw.selectedOrderRow < 0 {
		dialog.ShowInformation("Удаление заказа", "Строка не выбрана", mw.win)
		return
	}
	o, ok := mw.om.At(mw.selectedOrderRow)
	if !ok {
		return
	}
	dialog.ShowConfirm("Подтверждение",
		fmt.Sprintf("Удалить заказ кассира «%s» от %s %s?", o.CashierFIO, o.Date.String(), o.Time),
		func(ok bool) {
			if !ok {
				return
			}
			mw.om.Remove(mw.selectedOrderRow)
			mw.selectedOrderRow = -1
			mw.refresh()
		}, mw.win)
}

func (mw *MainWindow) onFindOrders() {
	if mw.cm.Size() == 0 {
		dialog.ShowInformation("Поиск заказов", "Нет кассиров для поиска", mw.win)
		return
	}
	fios, addrs := mw.uniqueFioAddr()
	fioSelect := widget.NewSelect(fios, nil)
	addrSelect := widget.NewSelect(addrs, nil)
	items := []*widget.FormItem{
		widget.NewFormItem("ФИО кассира", fioSelect),
		widget.NewFormItem("Адрес точки", addrSelect),
	}
	dialog.ShowForm("Найти заказы", "Найти", "Отмена", items, func(ok bool) {
		if !ok {
			return
		}
		orders, steps := mw.om.FindByCashierAddress(fioSelect.Selected, addrSelect.Selected)
		if len(orders) == 0 {
			dialog.ShowInformation("Поиск заказов",
				fmt.Sprintf("Заказы не найдены\nШагов поиска в БДП: %d", steps), mw.win)
			return
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Шагов поиска в БДП: %d\n\n", steps)
		for _, o := range orders {
			fmt.Fprintf(&sb, "%s  %s  %.2f руб.\n", o.Date.String(), o.Time, o.Sum)
		}
		label := widget.NewLabel(sb.String())
		d := dialog.NewCustom(fmt.Sprintf("Заказы: %d", len(orders)), "Закрыть", container.NewScroll(label), mw.win)
		d.Resize(fyne.NewSize(400, 400))
		d.Show()
	}, mw.win)
}

func (mw *MainWindow) onBuildReport() {
	if mw.cm.Size() == 0 {
		dialog.ShowInformation("Отчёт", "Нет доступных адресов", mw.win)
		return
	}
	_, addrs := mw.uniqueFioAddr()
	addrSelect := widget.NewSelect(addrs, nil)
	sumEntry := widget.NewEntry()
	sumEntry.SetText("300.00")
	dateEntry := widget.NewEntry()
	dateEntry.SetText("01.01.2025")

	items := []*widget.FormItem{
		widget.NewFormItem("Адрес точки", addrSelect),
		widget.NewFormItem("Сумма заказа", sumEntry),
		widget.NewFormItem("Дата отчёта", dateEntry),
	}
	dialog.ShowForm("Построить отчёт", "Построить", "Отмена", items, func(ok bool) {
		if !ok {
			return
		}
		sum, err := strconv.ParseFloat(strings.TrimSpace(sumEntry.Text), 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("report: invalid data(sum)"), mw.win)
			return
		}
		d, dok := date.Parse(strings.TrimSpace(dateEntry.Text))
		if !dok {
			dialog.ShowError(fmt.Errorf("report: invalid data(date)"), mw.win)
			return
		}
		rows := managers.BuildReport(mw.cm, mw.om, addrSelect.Selected, d, sum)
		if len(rows) == 0 {
			dialog.ShowInformation("Отчёт", "По заданным критериям данных не найдено", mw.win)
			return
		}
		var sb strings.Builder
		for _, r := range rows {
			fmt.Fprintf(&sb, "%s | %s | %.2f | %s %s | %.2f\n",
				r.CashierFIO, r.CashierAddress, r.CashierSalary, r.OrderDate, r.OrderTime, r.OrderSum)
		}
		label := widget.NewLabel(sb.String())
		d2 := dialog.NewCustom(fmt.Sprintf("Отчёт: %d строк", len(rows)), "Закрыть", container.NewScroll(label), mw.win)
		d2.Resize(fyne.NewSize(700, 400))
		d2.Show()
	}, mw.win)
}

func (mw *MainWindow) showLoadSummary(title string, loaded, skipped int, details []string) {
	text := fmt.Sprintf("Загружено: %d\nПропущено: %d", loaded, skipped)
	if len(details) == 0 {
		dialog.ShowInformation(title, text, mw.win)
		return
	}
	text += "\n\n" + strings.Join(details, "\n")
	label := widget.NewLabel(text)
	d := dialog.NewCustom(title, "Закрыть", container.NewScroll(label), mw.win)
	d.Resize(fyne.NewSize(520, 400))
	d.Show()
}

func (mw *MainWindow) confirmReplace(hasData bool, message string, proceed func()) {
	if !hasData {
		proceed()
		return
	}
	dialog.ShowConfirm("Подтверждение", message, func(ok bool) {
		if ok {
			proceed()
		}
	}, mw.win)
}

func (mw *MainWindow) onLoadCashiers() {
	mw.confirmReplace(mw.cm.Size() > 0, "Текущие кассиры и заказы будут удалены. Продолжить?", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil || r == nil {
				return
			}
			defer r.Close()
			mw.cm.Reset()
			mw.om.Reset()
			mw.selectedCashierRow = -1
			mw.selectedOrderRow = -1
			loaded, skipped, details, err := iodata.LoadCashiers(r, mw.cm)
			mw.refresh()
			if err != nil {
				dialog.ShowError(err, mw.win)
				return
			}
			mw.showLoadSummary("Загрузка кассиров", loaded, skipped, details)
		}, mw.win)
	})
}

func (mw *MainWindow) onLoadOrders() {
	if mw.cm.Size() == 0 {
		dialog.ShowInformation("Загрузить заказы", "Сначала загрузите кассиров", mw.win)
		return
	}
	mw.confirmReplace(mw.om.Size() > 0, "Текущие заказы будут удалены. Продолжить?", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil || r == nil {
				return
			}
			defer r.Close()
			mw.om.Reset()
			mw.selectedOrderRow = -1
			loaded, skipped, details, err := iodata.LoadOrders(r, mw.om, mw.cm)
			mw.refresh()
			if err != nil {
				dialog.ShowError(err, mw.win)
				return
			}
			mw.showLoadSummary("Загрузка заказов", loaded, skipped, details)
		}, mw.win)
	})
}

func (mw *MainWindow) saveTo(title string, count int, save func(io.Writer) error) {
	dialog.ShowFileSave(func(w fyne.URIWriteCloser, err error) {
		if err != nil || w == nil {
			return
		}
		err = save(w)
		if cerr := w.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			dialog.ShowError(err, mw.win)
			return
		}
		dialog.ShowInformation(title, fmt.Sprintf("Сохранено: %d", count), mw.win)
	}, mw.win)
}

func (mw *MainWindow) onSaveCashiers() {
	mw.saveTo("Сохранение кассиров", mw.cm.Size(), func(w io.Writer) error {
		return iodata.SaveCashiers(w, mw.cm)
	})
}

func (mw *MainWindow) onSaveOrders() {
	mw.saveTo("Сохранение заказов", mw.om.Size(), func(w io.Writer) error {
		return iodata.SaveOrders(w, mw.om)
	})
}
