package elements

import (
	"fmt"
	"time"

	"github.com/visualfc/goqt/ui"
	"github.com/waelbendhia/tariffs-app/types"
)

func newHistoryElement(app playtimeSeacher) (*ui.QGroupBox, func()) {
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
				if c != nil && c.IsWidgetType() {
					c.Delete()
				}
			}

			historyListLayout.AddWidget(newPTHeader())
			for _, pt := range pts {
				historyListLayout.AddWidget(newPTDisplay(pt))
			}
			historyListLayout.AddWidget(newPTTotal(pts))
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

	return rootBox, func() {
		withDates(func(min, max time.Time) (time.Time, time.Time) {
			updateSearch(min, max)
			return min, max
		})
	}
}

func newPTDisplay(pt types.Playtime) *ui.QWidget {
	var (
		shortFmt = "02/01/2006 15:04:05"
		b        = newPTRow(
			pt.Machine.Name,
			fmt.Sprintf("%d Millimes", pt.CalculatePrice()),
			pt.End.Sub(pt.Start).Truncate(time.Second).String(),
			pt.Start.Format(shortFmt),
			pt.End.Format(shortFmt),
		)
	)
	return b
}

func newPTHeader() *ui.QWidget {
	b := newPTRow(
		"Jeux",
		"Prix",
		"Durée",
		"Début",
		"Fin",
	)
	b.SetStyleSheet(`QLabel {
		font-weight: 900;
		color: rgba(0,0,0,0.6);
		font-size: 0.8em;
	}`)
	return b
}

func newPTTotal(pts []types.Playtime) *ui.QWidget {
	b := newPTRow(
		"Total:",
		fmt.Sprintf(
			"%d Millimes",
			func() int64 {
				var sum int64
				for _, pt := range pts {
					sum += pt.CalculatePrice()
				}
				return sum
			}(),
		),
		func() time.Duration {
			var sum time.Duration
			for _, pt := range pts {
				sum += pt.End.Sub(pt.Start)
			}
			return sum
		}().String(),
		"",
		"",
	)
	b.SetStyleSheet(`QLabel { font-weight: 900 }`)
	return b
}

func newPTRow(
	gameTxt string,
	priceTxt string,
	durationTxt string,
	startTxt string,
	endTxt string,
) *ui.QWidget {
	var (
		game     = newLabelWithText(gameTxt)
		price    = newLabelWithText(priceTxt)
		duration = newLabelWithText(durationTxt)
		start    = newLabelWithText(startTxt)
		end      = newLabelWithText(endTxt)
		oneFr    = ui.NewSizePolicyWithHorizontalVertical(
			ui.QSizePolicy_Minimum,
			ui.QSizePolicy_Minimum,
		)
		twoFr = ui.NewSizePolicyWithHorizontalVertical(
			ui.QSizePolicy_Minimum,
			ui.QSizePolicy_Minimum,
		)
		b, bLayout = newHBox()
	)
	oneFr.SetHorizontalStretch(1)
	twoFr.SetHorizontalStretch(2)
	game.SetSizePolicy(oneFr)
	price.SetSizePolicy(oneFr)
	price.SetAlignment(ui.Qt_AlignRight)
	duration.SetSizePolicy(oneFr)
	duration.SetAlignment(ui.Qt_AlignRight)
	start.SetSizePolicy(twoFr)
	start.SetAlignment(ui.Qt_AlignRight)
	end.SetSizePolicy(twoFr)
	end.SetAlignment(ui.Qt_AlignRight)
	bLayout.AddWidget(game)
	bLayout.AddWidget(price)
	bLayout.AddWidget(duration)
	bLayout.AddWidget(start)
	bLayout.AddWidget(end)
	return b
}
