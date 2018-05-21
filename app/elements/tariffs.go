package elements

import (
	"strconv"
	"time"

	"github.com/therecipe/qt/widgets"
	"github.com/waelbendhia/tariffs-app/types"
)

func newTariffElement(app tariffGetterSetter) *widgets.QGroupBox {
	var (
		// Holds the latest tariff
		tariffChan = make(chan *types.Tariff, 1)
		// withTariff provides a closure that accesses the value in tariffChan
		withTariff = func(f func(*types.Tariff) *types.Tariff) {
			t := <-tariffChan
			res := f(t)
			tariffChan <- res
		}
		// This will hold the user's input
		newTariff types.Tariff
		// Root Box will hold all the UI elements for tariff selection
		rootBox, rootLayout = newVGroupBoxWithTitle("Tariff")
		// These are the UI elements for setting the price per unit
		priceInputBox, priceInputLayout = newHBox()
		priceInput                      = widgets.NewQPlainTextEdit(nil)
		priceLabel                      = newLabelWithText("Prix en millime")
		// These are the UI elements for setting the unit for tarification
		unitInputBox, unitInputLayout = newHBox()
		unitInput                     = widgets.NewQPlainTextEdit(nil)
		unitSelection                 = newComboxBoxWithOptions(
			"Seconde",
			"Minute",
			"Heure",
		)
		unitLabel = newLabelWithText("Par")
		// These are the UI elements for buttons
		buttonBox, buttonLayout = newHBox()
		cancelButton            = widgets.NewQPushButton2("Anuller", nil)
		submitButton            = widgets.NewQPushButton2("Confirmer", nil)
		// toggleButton checks if newTariff is valid and disables submitButton
		// accordinlgy
		toggleButton = func() {
			var tariff *types.Tariff
			withTariff(func(t *types.Tariff) *types.Tariff {
				tariff = t
				return t
			})
			submitButton.SetEnabled(
				newTariff.PricePerUnit > 0 &&
					newTariff.UnitSize > 0 &&
					!newTariff.Equals(tariff),
			)
			cancelButton.SetEnabled(!newTariff.Equals(tariff))
		}
		// getUnit parses unit from unitInput
		getUnit = func() time.Duration {
			v, err := strconv.Atoi(unitInput.ToPlainText())
			var unit time.Duration = -1
			if err == nil && v > 0 {
				unit = time.Duration(v) *
					durationFromUnitSelectionInd(unitSelection.CurrentIndex())
			}
			return unit
		}
		// setTariffUI will update all our ui elements with given tariff
		setTariffUI = func(_ bool) {
			var tariff *types.Tariff
			withTariff(func(t *types.Tariff) *types.Tariff {
				tariff = t
				return t
			})
			if tariff != nil {
				priceInput.SetPlainText(strconv.Itoa(int(tariff.PricePerUnit)))
				unitInput.SetPlainText(strconv.Itoa(
					int(tariff.UnitSize) / int(
						durationFromUnitSelectionInd(
							unitSelectionIndFromDuration(tariff.UnitSize),
						)),
				))
				unitSelection.SetCurrentIndex(unitSelectionIndFromDuration(tariff.UnitSize))
			}
			toggleButton()
		}
		submitTariff = func(_ bool) {
			withTariff(func(t *types.Tariff) *types.Tariff {
				if newTariff.UnitSize > 0 &&
					newTariff.PricePerUnit > 0 &&
					!newTariff.Equals(t) {
					t2 := app.SetTariff(newTariff)
					t = &t2
				}
				return t
			})
			toggleButton()
		}
	)
	// First thing is we retrieve the latest tarifft
	tariffChan <- app.GetTariff()
	// Setup our elements based on the latest tariff
	withTariff(func(t *types.Tariff) *types.Tariff {
		if t != nil {
			newTariff.PricePerUnit = t.PricePerUnit
			newTariff.UnitSize = t.UnitSize
		}
		return t
	})
	setTariffUI(false)

	// Set up our the callbacks on inputs
	priceInput.ConnectTextChanged(func() {
		v, err := strconv.Atoi(priceInput.ToPlainText())
		if err != nil || v <= 0 {
			newTariff.PricePerUnit = -1
		} else {
			newTariff.PricePerUnit = int64(v)
		}
		toggleButton()
	})
	unitInput.ConnectTextChanged(func() {
		newTariff.UnitSize = getUnit()
		toggleButton()
	})

	unitSelection.ConnectCurrentIndexChanged(func(ind int) {
		newTariff.UnitSize = getUnit()
		toggleButton()
	})

	// Set SetTabChangesFocus on our inputs
	// priceInput.SetTabChangesFocus(true)
	// unitInput.SetTabChangesFocus(true)

	// Define button and input actions
	addOnEnterHandler(priceInput, func() { submitTariff(false) })
	addOnEnterHandler(unitInput, func() { submitTariff(false) })
	submitButton.ConnectClicked(submitTariff)
	cancelButton.ConnectClicked(setTariffUI)

	// Now we set up our ui elements
	buttonLayout.QLayout.AddWidget(cancelButton)
	buttonLayout.AddSpacerItem(
		widgets.NewQSpacerItem(
			0, 0,
			widgets.QSizePolicy__Expanding,
			widgets.QSizePolicy__Expanding,
		),
	)
	buttonLayout.QLayout.AddWidget(submitButton)
	buttonBox.SetMaximumHeight(inputHeight)

	priceInputLayout.QLayout.AddWidget(priceLabel)
	priceInputLayout.QLayout.AddWidget(priceInput)
	priceInputBox.SetMaximumHeight(inputHeight)

	unitInputLayout.QLayout.AddWidget(unitLabel)
	unitInputLayout.QLayout.AddWidget(unitInput)
	unitInputLayout.QLayout.AddWidget(unitSelection)
	unitInputBox.SetMaximumHeight(inputHeight)

	rootLayout.QLayout.AddWidget(priceInputBox)
	rootLayout.QLayout.AddWidget(unitInputBox)
	rootLayout.QLayout.AddWidget(buttonBox)

	rootBox.SetMaximumHeight(inputHeight * 4)

	rootBox.ConnectDestroyQObject(func() { close(tariffChan) })

	return rootBox
}

func unitSelectionIndFromDuration(dur time.Duration) int {
	switch {
	case dur%time.Hour == 0:
		return 2
	case dur%time.Minute == 0:
		return 1
	case dur%time.Second == 0:
		return 0
	default:
		return -1
	}
}

func durationFromUnitSelectionInd(ind int) time.Duration {
	switch ind {
	case 2:
		return time.Hour
	case 1:
		return time.Minute
	case 0:
		return time.Second
	default:
		return -1
	}
}
