package elements

import (
	"github.com/andlabs/ui"
)

func MainWindow(app tariffGetterSetter) *ui.Window {
	w := ui.NewWindow("Tariffs", 200, 100, true)
	t := tariffElement(app)
	w.SetChild(t)
	return w
}
