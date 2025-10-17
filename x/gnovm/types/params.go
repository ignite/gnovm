package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
)

var (
	defaultDepositCost int64 = 100
	defaultStorageCost int64 = 1
	moduleAccountAddr        = authtypes.NewModuleAddress(ModuleName)
)

// DefaultParams returns the default set of parameters.
func DefaultParams() Params {
	defaultVmParams := vm.DefaultParams()
	return Params{
		SysnamesPkgpath:     defaultVmParams.SysNamesPkgPath,
		ChainDomain:         defaultVmParams.ChainDomain,
		DefaultDeposit:      sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(defaultDepositCost)).String(),
		StoragePrice:        sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(defaultStorageCost)).String(),
		StorageFeeCollector: moduleAccountAddr,
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	return p.ToVmParams().Validate()
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

// VmParamsToParams converts the vm.Params to Params.
func VmParamsToParams(vmParams vm.Params) Params {
	return Params{
		SysnamesPkgpath:     vmParams.SysNamesPkgPath,
		ChainDomain:         vmParams.ChainDomain,
		DefaultDeposit:      vmParams.DefaultDeposit,
		StoragePrice:        vmParams.StoragePrice,
		StorageFeeCollector: vmParams.StorageFeeCollector[:],
	}
}
