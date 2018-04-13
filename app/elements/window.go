package elements

import (
	"github.com/visualfc/goqt/ui"
)

// MainWindow creates the app's main window
func MainWindow(app app) *ui.QMainWindow {

	var (
		w                  = ui.NewMainWindow()
		tab                = ui.NewTabWidget()
		suivi, suiviLayout = newHBox()
		// firstTab  = ui.NewHorizontalBox()
		tariff = newTariffElement(app)
		// seperator = ui.NewHorizontalSeparator()
		machines = newMachinesElement(app)
	)
	w.SetWindowTitle("Tariffs")

	suiviLayout.AddWidget(machines)
	suiviLayout.AddWidget(tariff)

	tab.AddTabWithWidgetString(suivi, "Suivie")

	w.SetCentralWidget(tab)

	// firstTab.Append(machines, true)
	// firstTab.Append(seperator, false)
	// firstTab.Append(tariff, false)

	// firstTab.SetPadded(true)

	// tab.Append("Tariff", firstTab)
	// tab.SetMargined(0, true)
	// w.SetChild(tab)
	return w
}
