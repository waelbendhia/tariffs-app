package elements

import (
	"time"

	"github.com/visualfc/goqt/ui"
)

const inputHeight int32 = 52

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

func newVGroupBox() (*ui.QGroupBox, *ui.QVBoxLayout) {
	b, l := ui.NewGroupBox(), ui.NewVBoxLayout()
	b.SetLayout(l)
	return b, l
}

func newVGroupBoxWithTitle(title string) (*ui.QGroupBox, *ui.QVBoxLayout) {
	b, l := ui.NewGroupBoxWithTitleParent(title, nil), ui.NewVBoxLayout()
	b.SetLayout(l)
	return b, l
}

func newLabelWithText(label string) *ui.QLabel {
	return ui.NewLabelWithTextParentFlags(label, nil, ui.Qt_Widget)
}

func newComboxBoxWithOptions(options ...string) *ui.QComboBox {
	b := ui.NewComboBox()
	b.AddItems(options)
	return b
}

func newVScroller() (*ui.QScrollArea, *ui.QVBoxLayout) {
	scroller := ui.NewScrollArea()
	layout := ui.NewVBoxLayout()
	return scroller, layout
}

type submitOnEnter struct {
	ui.QObject
	this     *ui.QWidget
	onSubmit func()
}

func newSubmitOnEnterFilter(f func()) *submitOnEnter {
	return &submitOnEnter{onSubmit: f}
}

func (soe *submitOnEnter) OnEvent(obj *ui.QObject, event *ui.QEvent) bool {
	if event.Type() == ui.QEvent_KeyPress {
		keyEvent := ui.QKeyEvent{
			QInputEvent: ui.QInputEvent{
				QEvent: *event,
			},
		}
		if ui.Qt_Key(keyEvent.Key()) == ui.Qt_Key_Return && soe.onSubmit != nil {
			soe.onSubmit()
			return true
		}
	}
	return soe.QObject.Event(event)
}

func newDateEdit() *ui.QDateEdit {
	de := ui.NewDateEdit()
	de.SetDisplayFormat("dd/MM/yyyy")
	de.SetDate(ui.QDateCurrentDate())
	return de
}

func qtDateToTime(dt ui.QDateTime) time.Time {
	return time.Unix(dt.ToMSecsSinceEpoch()/1000, 0)
}
