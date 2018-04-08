package elements

import (
	"strconv"
	"time"

	"github.com/andlabs/ui"
	"github.com/pkg/errors"
	"github.com/waelbendhia/tariffs-app/types"
)

func tariffElement(app tariffGetterSetter) *ui.Box {
	var (
		defTariff                         = app.GetTariff()
		inputTariff, getTariff, setTariff = tariffInput()
		fullBox                           = ui.NewVerticalBox()
		submit                            = func() { app.SetTariff(getTariff()) }
		cancel                            = func() { setTariff(app.GetTariff()) }
		buttonBox                         = buttonBox(cancel, submit)
	)

	setTariff(defTariff)

	fullBox.Append(inputTariff, false)
	fullBox.Append(buttonBox, false)
	return fullBox
}

func tariffInput() (
	*ui.Box,
	func() (*types.Tariff, error),
	func(*types.Tariff),
) {
	var (
		box                       = ui.NewVerticalBox()
		price, getPrice, setPrice = tariffPriceInput()
		unit, getUnit, setUnit    = tariffUnitInput()
		getTariff                 = func() (*types.Tariff, error) {
			price, pErr := getPrice()
			if pErr != nil {
				return nil, errors.Wrap(pErr, "Invalid input for tariff")
			}
			unit, uErr := getUnit()
			if uErr != nil {
				return nil, errors.Wrap(uErr, "Invalid input for tariff")
			}
			return &types.Tariff{
				PricePerUnit: price,
				UnitSize:     unit,
			}, nil
		}
		setTariff = func(t *types.Tariff) {
			if t != nil {
				setPrice(t.PricePerUnit)
				setUnit(t.UnitSize)
			}
		}
	)

	box.Append(price, false)
	box.Append(unit, false)

	return box, getTariff, setTariff

}

func tariffPriceInput() (
	*ui.Box,
	func() (int64, error),
	func(int64),
) {
	var (
		box      = ui.NewHorizontalBox()
		input    = ui.NewEntry()
		setPrice = func(p int64) {
			if p != -1 {
				input.SetText(strconv.Itoa(int(p)))
			}
		}
		getPrice = func() (int64, error) {
			v, err := strconv.Atoi(input.Text())
			if err != nil {
				return 0, errors.Errorf("Invalid input for price: %s", input.Text())
			}
			return int64(v), nil
		}
	)
	box.Append(ui.NewLabel("Prix en millime"), true)
	box.Append(input, false)
	return box, getPrice, setPrice
}

func tariffUnitInput() (
	*ui.Box,
	func() (time.Duration, error),
	func(time.Duration),
) {
	var (
		box        = ui.NewHorizontalBox()
		input      = ui.NewEntry()
		typeSelect = func() *ui.Combobox {
			s := ui.NewCombobox()
			s.Append("Heure")
			s.Append("Minute")
			s.Append("Seconde")
			return s
		}()
		getUnit = func() (time.Duration, error) {
			v, err := strconv.Atoi(input.Text())
			if err != nil {
				return 0, errors.Errorf("Invalid input for time: %s", input.Text())
			}
			var dur time.Duration
			switch typeSelect.Selected() {
			case 0:
				dur = time.Duration(v) * time.Hour
			case 1:
				dur = time.Duration(v) * time.Minute
			case 2:
				dur = time.Duration(v) * time.Second
			default:
				return 0, errors.New("No time type selected")
			}
			return dur, nil
		}
		setUnit = func(u time.Duration) {
			if u != -1 {
				switch {
				case u >= time.Hour:
					typeSelect.SetSelected(0)
					input.SetText(strconv.Itoa(int(u / time.Hour)))
				case u >= time.Minute:
					typeSelect.SetSelected(1)
					input.SetText(strconv.Itoa(int(u / time.Minute)))
				default:
					typeSelect.SetSelected(2)
					input.SetText(strconv.Itoa(int(u)))
				}
			}
		}
	)

	box.Append(ui.NewLabel("Unite de tariffation"), true)
	box.Append(input, false)
	box.Append(typeSelect, false)

	return box, getUnit, setUnit
}
func buttonBox(cancel, submit func()) *ui.Box {
	var (
		buttonBox    = ui.NewHorizontalBox()
		cancelButton = ui.NewButton("Annuler")
		sep          = ui.NewLabel("")
		submitButton = ui.NewButton("Confirmer")
	)
	buttonBox.Append(cancelButton, false)
	buttonBox.Append(sep, true)
	buttonBox.Append(submitButton, false)

	cancelButton.OnClicked(func(_ *ui.Button) {
		cancel()
	})

	submitButton.OnClicked(func(_ *ui.Button) {
		submit()
	})

	return buttonBox
}
