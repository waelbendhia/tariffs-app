package elements

import "github.com/waelbendhia/tariffs-app/types"

type tariffGetterSetter interface {
	GetTariff() *types.Tariff
	SetTariff(types.Tariff) types.Tariff
}

type machineCRUDER interface {
	GetMachines() []types.Machine
	DeleteMachine(m types.Machine) types.Machine
	AddMachine(m types.Machine) types.Machine
	UpdateMachine(m types.Machine) types.Machine
}

type app interface {
	tariffGetterSetter
	machineCRUDER
}
