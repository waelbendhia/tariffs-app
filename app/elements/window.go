package elements

import (
	"github.com/andlabs/ui"
)

// MainWindow creates the app's main window
func MainWindow(app tariffGetterSetter) *ui.Window {
	w := ui.NewWindow("Tariffs", 200, 100, true)
	w.SetChild(newTariffElement(app))
	return w
}
