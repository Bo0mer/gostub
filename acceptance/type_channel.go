package acceptance

import "github.com/momchil-atanasov/gostub/acceptance/external/external_dup"

//go:generate gostub ChannelSupport

type ChannelSupport interface {
	Method(chan external.Address) chan external.Address
}
