package elements

import (
	"strconv"
	"time"

	"github.com/visualfc/goqt/ui"
	"github.com/waelbendhia/tariffs-app/types"
)

func newTariffElement(app tariffGetterSetter) *ui.QGroupBox {
	var (
		// Holds the latest tariff
		tariffChan = make(chan *types.Tariff, 1)
		// This will hold the user's input
		newTariff types.Tariff
		// Root Box will hold all the UI elements for tariff selection
		rootBox, rootLayout = newVGroupBoxWithTitle("Tariff")
		// These are the UI elements for setting the price per unit
		priceInputBox, priceInputLayout = newHBox()
		priceInput                      = ui.NewPlainTextEdit()
		priceLabel                      = newLabelWithText("Prix en millime")
		// These are the UI elements for setting the unit for tarification
		unitInputBox, unitInputLayout = newHBox()
		unitInput                     = ui.NewPlainTextEdit()
		unitSelection                 = newComboxBoxWithOptions(
			"Seconde",
			"Minute",
			"Heure",
		)
		unitLabel = newLabelWithText("Par")
		// These are the UI elements for the dialog buttons
		buttonBox, buttonLayout = newHBox()
		cancelButton            = ui.NewPushButtonWithTextParent("Anuller", nil)
		submitButton            = ui.NewPushButtonWithTextParent("Confirmer", nil)
		// toggleButton checks if newTariff is valid and disables submitButton
		// accordinlgy
		toggleButton = func() {
			tariff := <-tariffChan
			submitButton.SetEnabled(
				newTariff.PricePerUnit > 0 &&
					newTariff.UnitSize > 0 &&
					!newTariff.Equals(tariff),
			)
			cancelButton.SetEnabled(!newTariff.Equals(tariff))
			tariffChan <- tariff
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
		setTariffUI = func() {
			t := <-tariffChan
			if t != nil {
				priceInput.SetPlainText(strconv.Itoa(int(t.PricePerUnit)))
				unitInput.SetPlainText(strconv.Itoa(
					int(t.UnitSize) / int(
						durationFromUnitSelectionInd(
							unitSelectionIndFromDuration(t.UnitSize),
						)),
				))
				unitSelection.SetCurrentIndex(unitSelectionIndFromDuration(t.UnitSize))
			}
			tariffChan <- t
			toggleButton()
		}
		submitTariff = func() {
			tariff := <-tariffChan
			if newTariff.UnitSize > 0 &&
				newTariff.PricePerUnit > 0 &&
				!newTariff.Equals(tariff) {
				t := app.SetTariff(newTariff)
				tariff = &t
			}
			tariffChan <- tariff
			toggleButton()
		}
	)
	// Setup our elements based on the latest tariff

	tariff := app.GetTariff()
	if tariff != nil {
		newTariff.PricePerUnit = tariff.PricePerUnit
		newTariff.UnitSize = tariff.UnitSize
	}
	tariffChan <- tariff
	setTariffUI()

	// Set up our inputs to update the newTariff
	priceInput.OnTextChanged(func() {
		v, err := strconv.Atoi(priceInput.ToPlainText())
		if err != nil || v <= 0 {
			newTariff.PricePerUnit = -1
		} else {
			newTariff.PricePerUnit = int64(v)
		}
		toggleButton()
	})
	priceInput.SetTabChangesFocus(true)
	unitInput.OnTextChanged(func() {
		newTariff.UnitSize = getUnit()
		toggleButton()
	})
	unitInput.SetTabChangesFocus(true)

	unitSelection.OnCurrentIndexChangedWithIndex(func(ind int32) {
		newTariff.UnitSize = getUnit()
		toggleButton()
	})
	// Define button and input actions
	priceInput.InstallEventFilter(newSubmitOnEnterFilter(submitTariff))
	unitInput.InstallEventFilter(newSubmitOnEnterFilter(submitTariff))
	submitButton.OnClicked(submitTariff)
	cancelButton.OnClicked(setTariffUI)

	// Now we set up our ui elements
	buttonLayout.AddWidget(cancelButton)
	buttonLayout.AddSpacerItem(
		ui.NewSpacerItem(0, 0, ui.QSizePolicy_Expanding, ui.QSizePolicy_Expanding),
	)
	buttonLayout.AddWidget(submitButton)
	buttonBox.SetMaximumHeight(inputHeight)

	priceInputLayout.AddWidget(priceLabel)
	priceInputLayout.AddWidget(priceInput)
	priceInputBox.SetMaximumHeight(inputHeight)

	unitInputLayout.AddWidget(unitLabel)
	unitInputLayout.AddWidget(unitInput)
	unitInputLayout.AddWidget(unitSelection)
	unitInputBox.SetMaximumHeight(inputHeight)

	rootLayout.AddWidget(priceInputBox)
	rootLayout.AddWidget(unitInputBox)
	rootLayout.AddWidget(buttonBox)

	rootBox.SetMaximumHeight(224)

	rootBox.OnDestroyed(func() { close(tariffChan) })

	return rootBox
}

func unitSelectionIndFromDuration(dur time.Duration) int32 {
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

func durationFromUnitSelectionInd(ind int32) time.Duration {
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
