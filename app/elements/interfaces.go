package elements

import "github.com/waelbendhia/tariffs-app/types"

type tariffGetter interface {
	GetTariff() *types.Tariff
}

type tariffSetter interface {
	SetTariff(*types.Tariff, error)
}

type tariffGetterSetter interface {
	tariffGetter
	tariffSetter
}
