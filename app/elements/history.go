package elements

import (
	"log"
	"time"

	"github.com/visualfc/goqt/ui"
)

func newHistoryElement(app playtimeSeacher) *ui.QGroupBox {
	var (
		minDateChan, maxDateChan   = make(chan time.Time, 1), make(chan time.Time, 1)
		stopChan                   = make(chan int)
		rootBox, rootBoxLayout     = newVGroupBoxWithTitle("Historique:")
		searchBox, searchBoxLayout = newHBox()
		minDateLabel               = newLabelWithText("Date minimum:")
		minDateInput               = newDateEdit()
		maxDateLabel               = newLabelWithText("Date maximum:")
		maxDateInput               = newDateEdit()
		updateSearch               = func(minDate, maxDate time.Time) {
			pts := app.SearchPlaytimes(nil, &minDate, &maxDate)
			var sum int64
			for _, pt := range pts {
				sum += pt.CalculatePrice()
			}
			log.Println(sum)
		}
	)

	minDateInput.OnDateTimeChanged(func(dt *ui.QDateTime) {
		maxDateInput.SetMinimumDateTime(dt)
		minDateChan <- qtDateToTime(*dt)
	})
	maxDateInput.OnDateTimeChanged(func(dt *ui.QDateTime) {
		minDateInput.SetMaximumDateTime(dt)
		maxDateChan <- qtDateToTime(*dt.AddSecs(3600 * 24))
	})

	searchBox.SetMaximumHeight(inputHeight)
	searchBox.SetLayout(searchBoxLayout)
	searchBoxLayout.AddWidget(minDateLabel)
	searchBoxLayout.AddWidget(minDateInput)
	searchBoxLayout.AddWidget(maxDateLabel)
	searchBoxLayout.AddWidget(maxDateInput)

	go func() {
		var minDate, maxDate time.Time
		for {
			select {
			case <-stopChan:
				return
			case minDate = <-minDateChan:
				updateSearch(minDate, maxDate)
			case maxDate = <-maxDateChan:
				updateSearch(minDate, maxDate)
			}
		}
	}()

	app.SearchPlaytimes(nil, nil, nil)

	rootBoxLayout.AddWidget(searchBox)

	rootBox.SetLayout(rootBoxLayout)
	return rootBox
}
