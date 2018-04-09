package elements

import (
	"strconv"
	"time"

	"github.com/andlabs/ui"
	"github.com/waelbendhia/tariffs-app/types"
)

func newTariffElement(app tariffGetterSetter) *ui.Box {
	var (
		// Root Box will hold all the UI elements for tariff selection
		rootBox = ui.NewVerticalBox()
		// These are the UI elements for setting the price per unit
		priceInputBox = ui.NewHorizontalBox()
		priceLabel    = ui.NewLabel("Prix en millime")
		priceInput    = ui.NewEntry()
		// These are the UI elements for setting the unit for tarification
		unitInputBox  = ui.NewHorizontalBox()
		unitLabel     = ui.NewLabel("Par")
		unitInput     = ui.NewEntry()
		unitSelection = unitSelectionBox()
		// These are the UI elements for the dialog buttons
		buttonBox    = ui.NewHorizontalBox()
		cancelButton = ui.NewButton("Annuler")
		seperator    = ui.NewLabel("")
		submitButton = ui.NewButton("Confirmer")
		// Get the initial state for this UI element
		tariff = app.GetTariff()
		// This will hold the user's input
		newTariff types.Tariff
		// toggleButton checks if newTariff is valid and disables submitButton
		// accordinlgy
		toggleButton = func() {
			if newTariff.PricePerUnit > 0 &&
				newTariff.UnitSize > 0 &&
				!newTariff.Equals(tariff) {
				submitButton.Enable()
			} else {
				submitButton.Disable()
			}
			if newTariff.Equals(tariff) {
				cancelButton.Disable()
			} else {
				cancelButton.Enable()
			}
		}
		// getUnit parses unit from unitInput
		getUnit = func() time.Duration {
			v, err := strconv.Atoi(unitInput.Text())
			var unit time.Duration = -1
			if err == nil && v > 0 {
				unit = time.Duration(v) *
					durationFromUnitSelectionInd(unitSelection.Selected())
			}
			return unit
		}
		// setTariffUI will update all our ui elements with given tariff
		setTariffUI = func(t *types.Tariff) {
			if t != nil {
				priceInput.SetText(strconv.Itoa(int(t.PricePerUnit)))
				unitInput.SetText(strconv.Itoa(
					int(t.UnitSize) / int(
						durationFromUnitSelectionInd(
							unitSelectionIndFromDuration(t.UnitSize),
						)),
				))
				unitSelection.SetSelected(unitSelectionIndFromDuration(t.UnitSize))
			}
			toggleButton()
		}
	)
	// Setup our elements based on the latest tariff
	if tariff != nil {
		newTariff.PricePerUnit = tariff.PricePerUnit
		newTariff.UnitSize = tariff.UnitSize
	}
	setTariffUI(tariff)

	// Set up our inputs to update the newTariff
	priceInput.OnChanged(func(e *ui.Entry) {
		v, err := strconv.Atoi(e.Text())
		if err != nil || v <= 0 {
			newTariff.PricePerUnit = -1
		} else {
			newTariff.PricePerUnit = int64(v)
		}
		toggleButton()
	})
	unitInput.OnChanged(func(_ *ui.Entry) {
		newTariff.UnitSize = getUnit()
		toggleButton()
	})
	unitSelection.OnSelected(func(_ *ui.Combobox) {
		newTariff.UnitSize = getUnit()
		toggleButton()
	})

	// Define button actions
	submitButton.OnClicked(func(_ *ui.Button) {
		t := app.SetTariff(newTariff)
		tariff = &t
		toggleButton()
	})
	cancelButton.OnClicked(func(_ *ui.Button) {
		setTariffUI(tariff)
	})

	// Now we set up our ui elements
	buttonBox.Append(cancelButton, false)
	buttonBox.Append(seperator, true)
	buttonBox.Append(submitButton, true)

	buttonBox.SetPadded(true)

	priceInputBox.Append(priceLabel, true)
	priceInputBox.Append(priceInput, false)

	priceInputBox.SetPadded(true)

	unitInputBox.Append(unitLabel, true)
	unitInputBox.Append(unitInput, false)
	unitInputBox.Append(unitSelection, false)

	unitInputBox.SetPadded(true)

	rootBox.Append(priceInputBox, false)
	rootBox.Append(unitInputBox, false)
	rootBox.Append(buttonBox, false)

	rootBox.SetPadded(true)

	return rootBox
}

func unitSelectionBox() *ui.Combobox {
	c := ui.NewCombobox()
	c.Append("Seconde")
	c.Append("Minute")
	c.Append("Heure")
	return c
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
