package elements

import (
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func newHBox() (*widgets.QWidget, *widgets.QHBoxLayout) {
	b, l := widgets.NewQWidget(nil, 0), widgets.NewQHBoxLayout()
	b.SetLayout(l)
	return b, l
}

func newVBox() (*widgets.QWidget, *widgets.QVBoxLayout) {
	b, l := widgets.NewQWidget(nil, 0), widgets.NewQVBoxLayout()
	b.SetLayout(l)
	return b, l
}

func newVGroupBox() (*widgets.QGroupBox, *widgets.QVBoxLayout) {
	b, l := widgets.NewQGroupBox(nil), widgets.NewQVBoxLayout()
	b.SetLayout(l)
	return b, l
}

func newVGroupBoxWithTitle(title string) (*widgets.QGroupBox, *widgets.QVBoxLayout) {
	b, l := widgets.NewQGroupBox2(title, nil), widgets.NewQVBoxLayout()
	b.SetLayout(l)
	return b, l
}

func newLabelWithText(label string) *widgets.QLabel {
	return widgets.NewQLabel2(label, nil, core.Qt__Widget)
}

func newComboxBoxWithOptions(options ...string) *widgets.QComboBox {
	b := widgets.NewQComboBox(nil)
	b.AddItems(options)
	return b
}

func newVScroller() (*widgets.QScrollArea, *widgets.QVBoxLayout) {
	scroller := widgets.NewQScrollArea(nil)
	layout := widgets.NewQVBoxLayout()
	return scroller, layout
}

type submitOnEnter struct {
	core.QObject
	this     *widgets.QWidget
	onSubmit func()
}

func addOnEnterHandler(watched *widgets.QPlainTextEdit, f func()) {
	watched.ConnectKeyPressEvent(
		func(event *gui.QKeyEvent) {
			if core.Qt__Key(event.Key()) == core.Qt__Key_Return && f != nil {
				f()
			} else {
				watched.KeyPressEventDefault(event)
			}
		})
}

func newDateEdit() *widgets.QDateEdit {
	de := widgets.NewQDateEdit(nil)
	de.SetDisplayFormat("dd/MM/yyyy")
	de.SetDate(core.QDate_CurrentDate())
	return de
}

func qtDateToTime(dt core.QDateTime) time.Time {
	return time.Unix(dt.ToMSecsSinceEpoch()/1000, 0)
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func newButton(text string) *widgets.QPushButton {
	return widgets.NewQPushButton2(text, nil)
}
