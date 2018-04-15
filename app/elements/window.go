package elements

import (
	"github.com/visualfc/goqt/ui"
)

// MainWindow creates the app's main window
func MainWindow(a app) *ui.QMainWindow {
	// Set input heights based on font metrics
	f := ui.NewFont()
	f.SetPointSize(18)
	inputHeight = ui.NewFontMetrics(f).Height() * 3 / 2

	var (
		aw                 = newAppWrapper(a)
		w                  = ui.NewMainWindow()
		suivi, suiviLayout = newHBox()
		rightPane          = ui.NewVBoxLayout()
		history            *ui.QGroupBox
		machines           *ui.QGroupBox
		tariff             *ui.QGroupBox
	)
	history, aw.ptUpdateFn = newHistoryElement(aw)
	machines, aw.trfUpdateFn = newMachinesElement(aw)
	tariff = newTariffElement(aw)
	w.SetWindowTitle("Tariffs")

	suiviLayout.AddWidget(machines)
	rightPane.AddWidget(tariff)
	rightPane.AddWidget(history)
	suiviLayout.AddLayout(rightPane)

	w.SetCentralWidget(suivi)

	return w
}
