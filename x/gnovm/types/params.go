package types

import "github.com/gnolang/gno/gno.land/pkg/sdk/vm"

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	defaultVmParams := vm.DefaultParams()
	return Params{
		SysnamesPkgpath: defaultVmParams.SysNamesPkgPath,
		ChainDomain:     defaultVmParams.ChainDomain,
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	vmParams := vm.Params{
		SysNamesPkgPath: p.SysnamesPkgpath,
		ChainDomain:     p.ChainDomain,
	}

	return vmParams.Validate()
}
