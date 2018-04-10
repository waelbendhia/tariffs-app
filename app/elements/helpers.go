package elements

import "github.com/visualfc/goqt/ui"

func newHBox() (*ui.QWidget, *ui.QHBoxLayout) {
	b, l := ui.NewWidget(), ui.NewHBoxLayout()
	b.SetLayout(l)
	return b, l
}

func newVBox() (*ui.QWidget, *ui.QVBoxLayout) {
	b, l := ui.NewWidget(), ui.NewVBoxLayout()
	b.SetLayout(l)
	return b, l
}
