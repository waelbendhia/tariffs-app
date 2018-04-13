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
		tariff             = newTariffElement(app)
		machines           = newMachinesElement(app)
	)
	w.SetWindowTitle("Tariffs")

	suiviLayout.AddWidget(machines)
	suiviLayout.AddWidget(tariff)

	tab.AddTabWithWidgetString(suivi, "Suivie")

	w.SetCentralWidget(tab)

	return w
}
