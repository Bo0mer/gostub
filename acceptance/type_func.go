package acceptance

import "github.com/momchil-atanasov/gostub/acceptance/external/external_dup"

//go:generate gostub FuncSupport

type FuncSupport interface {
	Method(func(external.Address) external.Address) func(external.Address) external.Address
}
