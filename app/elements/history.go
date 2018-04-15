package elements

import (
	"fmt"
	"time"

	"github.com/visualfc/goqt/ui"
)

func newHistoryElement(app playtimeSeacher) *ui.QGroupBox {
	var (
		minDateChan, maxDateChan = make(chan time.Time, 1), make(chan time.Time, 1)
		withDates                = func(
			f func(time.Time, time.Time) (time.Time, time.Time),
		) {
			minDate, maxDate := <-minDateChan, <-maxDateChan
			minDate, maxDate = f(minDate, maxDate)
			minDateChan <- minDate
			maxDateChan <- maxDate
		}
		rootBox, rootBoxLayout         = newVGroupBoxWithTitle("Historique:")
		searchBox, searchBoxLayout     = newHBox()
		minDateLabel                   = newLabelWithText("Date minimum:")
		minDateInput                   = newDateEdit()
		maxDateLabel                   = newLabelWithText("Date maximum:")
		maxDateInput                   = newDateEdit()
		historyScroller                = ui.NewScrollArea()
		historyList, historyListLayout = newVBox()
		updateSearch                   = func(minDate, maxDate time.Time) {
			pts := app.SearchPlaytimes(nil, &minDate, &maxDate)

			for _, c := range historyList.Children() {
				if c.IsWidgetType() {
					c.Delete()
				}
			}

			for _, pt := range pts {
				historyListLayout.AddWidget(
					newLabelWithText(
						fmt.Sprintf("%s %d %v %v", pt.Machine.Name, pt.CalculatePrice(), pt.Start, pt.End),
					),
				)
			}
		}
	)

	minDateInput.OnDateTimeChanged(func(dt *ui.QDateTime) {
		withDates(func(minDate, maxDate time.Time) (time.Time, time.Time) {
			minDate = qtDateToTime(*dt)
			updateSearch(minDate, maxDate)
			return minDate, maxDate
		})
		maxDateInput.SetMinimumDateTime(dt)
	})
	maxDateInput.OnDateTimeChanged(func(dt *ui.QDateTime) {
		withDates(func(minDate, maxDate time.Time) (time.Time, time.Time) {
			maxDate = qtDateToTime(*dt.AddSecs(3600 * 24))
			updateSearch(minDate, maxDate)
			return minDate, maxDate
		})
		minDateInput.SetMaximumDateTime(dt)
	})

	searchBox.SetMaximumHeight(inputHeight)
	searchBox.SetLayout(searchBoxLayout)
	searchBoxLayout.AddWidget(minDateLabel)
	searchBoxLayout.AddWidget(minDateInput)
	searchBoxLayout.AddWidget(maxDateLabel)
	searchBoxLayout.AddWidget(maxDateInput)

	historyList.SetLayout(historyListLayout)
	historyListLayout.SetAlignment(ui.Qt_AlignTop)

	historyScroller.SetWidget(historyList)
	historyScroller.SetWidgetResizable(true)

	rootBoxLayout.AddWidget(searchBox)
	rootBoxLayout.AddWidget(historyScroller)

	rootBox.SetLayout(rootBoxLayout)
	minDateChan <- truncateToDay(time.Now())
	maxDateChan <- truncateToDay(time.Now().Add(24 * time.Hour))
	withDates(func(min, max time.Time) (time.Time, time.Time) {
		updateSearch(min, max)
		return min, max
	})
	return rootBox
}
