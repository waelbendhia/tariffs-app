package elements

import (
	"time"

	"github.com/waelbendhia/tariffs-app/types"
)

type tariffGetterSetter interface {
	GetTariff() *types.Tariff
	SetTariff(types.Tariff) types.Tariff
}

type timer interface {
	GetOpenPlayTime(machineID int64) *types.Playtime
	Start(machineID int64) types.Playtime
	End(ID int64) types.Playtime
}

type machineCRUDER interface {
	GetMachines() []types.Machine
	DeleteMachine(m types.Machine) types.Machine
	AddMachine(m types.Machine) types.Machine
}

type playtimeSeacher interface {
	SearchPlaytimes(
		machineID *int64,
		minDate *time.Time,
		maxDate *time.Time,
	) []types.Playtime
}

type machineCRUDERTimer interface {
	timer
	machineCRUDER
}

type app interface {
	tariffGetterSetter
	machineCRUDERTimer
	playtimeSeacher
}

type appWrapper struct {
	app
	trfUpdateFn func(*types.Tariff)
	ptUpdateFn  func()
}

func newAppWrapper(a app) *appWrapper {
	return &appWrapper{app: a}
}

func (aw *appWrapper) GetTariff() *types.Tariff {
	t := aw.app.GetTariff()
	if aw.trfUpdateFn != nil {
		aw.trfUpdateFn(t)
	}
	return t
}

func (aw *appWrapper) SetTariff(t types.Tariff) types.Tariff {
	t = aw.app.SetTariff(t)
	if aw.trfUpdateFn != nil {
		aw.trfUpdateFn(&t)
	}
	return t
}

func (aw *appWrapper) End(ID int64) types.Playtime {
	res := aw.app.End(ID)
	if aw.ptUpdateFn != nil {
		aw.ptUpdateFn()
	}
	return res
}
