package elements

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// MainWindow creates the app's main window
func MainWindow(a app) *widgets.QMainWindow {
	// Set input heights based on font metrics
	f := gui.NewQFont()
	f.SetPointSize(18)
	inputHeight = gui.NewQFontMetrics(f).Height() * 5 / 3

	var (
		aw                 = newAppWrapper(a)
		w                  = widgets.NewQMainWindow(nil, 0)
		suivi, suiviLayout = newHBox()
		rightPane          = widgets.NewQVBoxLayout()
		history            *widgets.QGroupBox
		machines           *widgets.QGroupBox
		tariff             *widgets.QGroupBox
	)
	history, aw.ptUpdateFn = newHistoryElement(aw)
	machines, aw.trfUpdateFn = newMachinesElement(aw)
	tariff = newTariffElement(aw)
	w.SetWindowTitle("Tariffs")

	suiviLayout.QLayout.AddWidget(machines)
	rightPane.QLayout.AddWidget(tariff)
	rightPane.QLayout.AddWidget(history)
	suiviLayout.AddLayout(rightPane, 0)

	w.SetCentralWidget(suivi)

	return w
}
