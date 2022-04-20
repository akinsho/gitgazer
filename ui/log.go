package ui

import (
	"akinsho/gitgazer/app"
	"fmt"
	"time"

	"github.com/rivo/tview"
)

type LogWidget struct {
	component *tview.TextView
}

func (d *LogWidget) Write(str string) {
	previous := d.Read()
	now := time.Now().Format("02-01-2006 15:04:05")
	output := fmt.Sprintf("[-:-:b]%s[::-]: %s", now, str)
	if previous != "" {
		output = previous + output
	}
	d.component.SetText(output)
}

func (d *LogWidget) Read() string {
	return d.component.GetText(false)
}

func (d *LogWidget) Refresh() (err error) {
	return
}

func (d *LogWidget) Component() tview.Primitive {
	return d.component
}

func (d *LogWidget) IsEmpty() bool {
	return d.component.GetText(false) == ""
}

func logWidget(ctx *app.Context) *LogWidget {
	debug := tview.NewTextView()
	debug.SetDynamicColors(true).SetBorder(true).SetTitle("Debug")
	return &LogWidget{debug}
}
