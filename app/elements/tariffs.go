package elements

import (
	"strconv"
	"time"

	"github.com/andlabs/ui"
	"github.com/pkg/errors"
	"github.com/waelbendhia/tariffs-app/types"
)

// TariffInput creates a new tariff input setting values to t if t is not nul
// and using fOnSubmit to handle user input
func TariffInput(t *types.Tariff, fOnSubmit func(*types.Tariff, error)) *ui.Box {
	button := submitButton()
	tInput, tGet := tariffInput(t)
	button.OnClicked(func(b *ui.Button) {
		t, err := tGet()
		fOnSubmit(t, err)
	})
	fullBox := ui.NewVerticalBox()
	fullBox.Append(tInput, false)
	fullBox.Append(button, false)
	return fullBox
}

func tariffInput(t *types.Tariff) (*ui.Box, func() (*types.Tariff, error)) {
	var (
		box                    = ui.NewVerticalBox()
		defPrice int64         = -1
		defUnit  time.Duration = -1
	)
	if t != nil {
		defPrice, defUnit = t.PricePerUnit, t.UnitSize
	}
	price, priceF := tariffPriceInput(defPrice)
	unit, unitF := tariffUnitInput(defUnit)

	box.Append(price, false)
	box.Append(unit, false)

	return box, func() (*types.Tariff, error) {
		price, pErr := priceF()
		if pErr != nil {
			return nil, errors.Wrap(pErr, "Invalid input for tariff")
		}
		unit, uErr := unitF()
		if uErr != nil {
			return nil, errors.Wrap(uErr, "Invalid input for tariff")
		}

		return &types.Tariff{
			PricePerUnit: price,
			UnitSize:     unit,
		}, nil
	}
}

func tariffPriceInput(defPrice int64) (*ui.Box, func() (int64, error)) {
	box := ui.NewHorizontalBox()
	box.Append(ui.NewLabel("Prix en millime"), false)
	input := ui.NewEntry()
	if defPrice != -1 {
		input.SetText(strconv.Itoa(int(defPrice)))
	}
	box.Append(input, false)
	return box, func() (int64, error) {
		v, err := strconv.Atoi(input.Text())
		if err != nil {
			return 0, errors.Errorf("Invalid input for price: %s", input.Text())
		}
		return int64(v), nil
	}
}

func tariffUnitInput(defUnit time.Duration) (*ui.Box, func() (time.Duration, error)) {
	box := ui.NewHorizontalBox()
	box.Append(ui.NewLabel("Unite de tariffation"), false)
	input := ui.NewEntry()

	box.Append(input, false)
	typeSelect := ui.NewCombobox()
	typeSelect.Append("Heure")
	typeSelect.Append("Minute")
	typeSelect.Append("Seconde")
	typeSelect.Selected()

	if defUnit != -1 {
		switch {
		case defUnit >= time.Hour:
			typeSelect.SetSelected(0)
			input.SetText(strconv.Itoa(int(defUnit / time.Hour)))
		case defUnit >= time.Minute:
			typeSelect.SetSelected(1)
			input.SetText(strconv.Itoa(int(defUnit / time.Minute)))
		default:
			typeSelect.SetSelected(2)
			input.SetText(strconv.Itoa(int(defUnit)))
		}
	}
	box.Append(typeSelect, false)
	return box,
		func() (time.Duration, error) {
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
}

func submitButton() *ui.Button {
	return ui.NewButton("Confirmer")
}
