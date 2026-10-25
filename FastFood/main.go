package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fastfood/ui"
)

func main() {
	a := app.New()
	win := a.NewWindow("Точка общепита")
	win.Resize(fyne.NewSize(400, 150))

	hashEntry := widget.NewEntry()
	hashEntry.SetText("17")

	dialog.ShowForm("Размер хеш-таблицы", "ОК", "Отмена",
		[]*widget.FormItem{widget.NewFormItem("Размер (15-20)", hashEntry)},
		func(ok bool) {
			if !ok {
				win.Close()
				return
			}
			size, err := strconv.Atoi(strings.TrimSpace(hashEntry.Text))
			invalid := err != nil || size < 15 || size > 20
			if invalid {
				size = 17
			}
			ui.NewMainWindow(win, size)
			if invalid {
				dialog.ShowInformation("Размер хеш-таблицы",
					fmt.Sprintf("Некорректный размер «%s», используется %d", hashEntry.Text, size), win)
			}
		}, win)

	win.ShowAndRun()
}
