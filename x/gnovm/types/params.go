package types

import (
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
)

// DefaultParams returns the default set of parameters.
func DefaultParams() Params {
	defaultVmParams := vm.DefaultParams()
	return Params{
		SysnamesPkgpath:     defaultVmParams.SysNamesPkgPath,
		ChainDomain:         defaultVmParams.ChainDomain,
		DefaultDeposit:      defaultVmParams.DefaultDeposit,
		StoragePrice:        defaultVmParams.StoragePrice,
		StorageFeeCollector: defaultVmParams.StorageFeeCollector[:],
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	vmParams := vm.Params{
		SysNamesPkgPath:     p.SysnamesPkgpath,
		ChainDomain:         p.ChainDomain,
		DefaultDeposit:      p.DefaultDeposit,
		StoragePrice:        p.StoragePrice,
		StorageFeeCollector: ToCryptoAddress(p.StorageFeeCollector),
	}

	return vmParams.Validate()
}

// ToVmParams converts the Params to vm.Params.
func (p Params) ToVmParams() vm.Params {
	vmParams := vm.Params{
		SysNamesPkgPath:     p.SysnamesPkgpath,
		ChainDomain:         p.ChainDomain,
		DefaultDeposit:      p.DefaultDeposit,
		StoragePrice:        p.StoragePrice,
		StorageFeeCollector: ToCryptoAddress(p.StorageFeeCollector),
	}

	return vmParams
}
