package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type tableRow struct {
	widget.BaseWidget
	labels      []*widget.Label
	row         int
	onTap       func(row int)
	onSecondary func(row int, pos fyne.Position)
}

func newTableRow(cols int, onTap func(int), onSecondary func(int, fyne.Position)) *tableRow {
	r := &tableRow{onTap: onTap, onSecondary: onSecondary}
	for i := 0; i < cols; i++ {
		r.labels = append(r.labels, widget.NewLabel(""))
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *tableRow) CreateRenderer() fyne.WidgetRenderer {
	objs := make([]fyne.CanvasObject, len(r.labels))
	for i, l := range r.labels {
		objs[i] = l
	}
	return widget.NewSimpleRenderer(container.NewGridWithColumns(len(objs), objs...))
}

func (r *tableRow) SetHeader(texts ...string) {
	for i, t := range texts {
		r.labels[i].TextStyle = fyne.TextStyle{Bold: true}
		r.labels[i].SetText(t)
		r.labels[i].Refresh()
	}
}

func (r *tableRow) SetTexts(selected bool, texts ...string) {
	marker := "   "
	if selected {
		marker = "▶  "
	}
	for i, t := range texts {
		if i == 0 {
			r.labels[i].SetText(marker + t)
		} else {
			r.labels[i].SetText(t)
		}
	}
}

func (r *tableRow) Tapped(_ *fyne.PointEvent) {
	if r.onTap != nil {
		r.onTap(r.row)
	}
}

func (r *tableRow) TappedSecondary(e *fyne.PointEvent) {
	if r.onTap != nil {
		r.onTap(r.row)
	}
	if r.onSecondary != nil {
		r.onSecondary(r.row, e.AbsolutePosition)
	}
}
