package elements

import (
	"github.com/visualfc/goqt/ui"
)

// MainWindow creates the app's main window
func MainWindow(a app) *ui.QMainWindow {
	var (
		aw                 = newAppWrapper(a)
		w                  = ui.NewMainWindow()
		suivi, suiviLayout = newHBox()
		rightPane          = ui.NewVBoxLayout()
		tariff             = newTariffElement(aw)
		machines           = newMachinesElement(aw, aw.tariffChan)
		history            = newHistoryElement(a)
	)
	w.SetWindowTitle("Tariffs")

	suiviLayout.AddWidget(machines)
	rightPane.AddWidget(tariff)
	rightPane.AddWidget(history)
	suiviLayout.AddLayout(rightPane)

	w.SetCentralWidget(suivi)
	w.OnDestroyed(func() { close(aw.tariffChan) })
	return w
}
