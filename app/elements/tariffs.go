package elements

import (
	"strconv"
	"time"

	"github.com/visualfc/goqt/ui"
	"github.com/waelbendhia/tariffs-app/types"
)

func newTariffElement(app tariffGetterSetter) *ui.QWidget {
	var (
		// Root Box will hold all the UI elements for tariff selection
		rootBox, rootLayout = newVBox()
		// These are the UI elements for setting the price per unit
		priceInputBox, priceInputLayout = newHBox()
		priceInput                      = ui.NewPlainTextEdit()
		priceLabel                      = ui.NewLabelWithTextParentFlags(
			"Prix en millime",
			nil,
			ui.Qt_Widget,
		)
		// These are the UI elements for setting the unit for tarification
		unitInputBox, unitInputLayout = newHBox()
		unitInput                     = ui.NewPlainTextEdit()
		unitSelection                 = func() *ui.QComboBox {
			b := ui.NewComboBox()
			b.AddItem("Seconde")
			b.AddItem("Minute")
			b.AddItem("Heure")
			return b
		}()
		unitLabel = ui.NewLabelWithTextParentFlags(
			"Par",
			nil,
			ui.Qt_Widget,
		)
		// These are the UI elements for the dialog buttons
		buttonBox, buttonLayout = newHBox()
		cancelButton            = ui.NewPushButtonWithTextParent("Anuller", nil)
		submitButton            = ui.NewPushButtonWithTextParent("Confirmer", nil)
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
				submitButton.SetEnabled(true)
			} else {
				submitButton.SetEnabled(false)
			}
			if newTariff.Equals(tariff) {
				cancelButton.SetEnabled(false)
			} else {
				cancelButton.SetEnabled(true)
			}
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
		setTariffUI = func(t *types.Tariff) {
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
	priceInput.OnTextChanged(func() {
		v, err := strconv.Atoi(priceInput.ToPlainText())
		if err != nil || v <= 0 {
			newTariff.PricePerUnit = -1
		} else {
			newTariff.PricePerUnit = int64(v)
		}
		toggleButton()
	})
	unitInput.OnTextChanged(func() {
		newTariff.UnitSize = getUnit()
		toggleButton()
	})
	unitSelection.OnCurrentIndexChangedWithIndex(func(ind int32) {
		newTariff.UnitSize = getUnit()
		toggleButton()
	})

	// Define button actions
	submitButton.OnClicked(func() {
		t := app.SetTariff(newTariff)
		tariff = &t
		toggleButton()
	})
	cancelButton.OnClicked(func() {
		setTariffUI(tariff)
	})

	// Now we set up our ui elements
	buttonLayout.AddWidget(cancelButton)
	buttonLayout.AddSpacerItem(
		ui.NewSpacerItem(0, 0, ui.QSizePolicy_Expanding, ui.QSizePolicy_Expanding),
	)
	buttonLayout.AddWidget(submitButton)
	buttonBox.SetMaximumHeight(48)

	priceInputLayout.AddWidget(priceLabel)
	priceInputLayout.AddWidget(priceInput)
	priceInputBox.SetMaximumHeight(48)

	unitInputLayout.AddWidget(unitLabel)
	unitInputLayout.AddWidget(unitInput)
	unitInputLayout.AddWidget(unitSelection)
	unitInputBox.SetMaximumHeight(48)

	rootLayout.AddWidget(priceInputBox)
	rootLayout.AddWidget(unitInputBox)
	rootLayout.AddWidget(buttonBox)

	rootBox.SetMaximumHeight(224)

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
