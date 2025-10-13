package gnovm

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Logger       log.Logger
	Config       *types.Module
	StoreKey     *storetypes.KVStoreKey
	MemStoreKey  *storetypes.MemoryStoreKey
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper types.AuthKeeper
	BankKeeper types.BankKeeper
}

type ModuleOutputs struct {
	depinject.Out

	GnoVMKeeper keeper.Keeper
	Module      appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.Logger,
		in.StoreKey,
		in.MemStoreKey,
		in.Cdc,
		in.AddressCodec,
		authority,
		in.AuthKeeper,
		in.BankKeeper,
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{GnoVMKeeper: k, Module: m}
}
