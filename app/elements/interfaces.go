package elements

import "github.com/waelbendhia/tariffs-app/types"

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
	UpdateMachine(m types.Machine) types.Machine
}

type machineCRUDERTimer interface {
	timer
	machineCRUDER
}

type app interface {
	tariffGetterSetter
	machineCRUDERTimer
}
