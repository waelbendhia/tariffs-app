package elements

import (
	"github.com/andlabs/ui"
)

// MainWindow creates the app's main window
func MainWindow(app app) *ui.Window {
	var (
		w         = ui.NewWindow("Tariffs", 200, 100, true)
		tab       = ui.NewTab()
		firstTab  = ui.NewHorizontalBox()
		tariff    = newTariffElement(app)
		seperator = ui.NewHorizontalSeparator()
		machines  = newMachinesElement(app)
	)

	firstTab.Append(machines, true)
	firstTab.Append(seperator, false)
	firstTab.Append(tariff, false)

	firstTab.SetPadded(true)

	tab.Append("Tariff", firstTab)
	tab.SetMargined(0, true)
	w.SetChild(tab)
	return w
}
