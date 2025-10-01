package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/ignite/gnovm/x/gnovm/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.VMKeeper.InitGenesis(
		types.GnoContextFromSDKContext(sdkCtx, k.logger),
		vm.GenesisState{ // todo: module params from the module itself and from the vmkeeper must stay in sync
			Params: genState.Params.ToVmParams(),
		},
	)

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	vmGenState := k.VMKeeper.ExportGenesis(types.GnoContextFromSDKContext(sdkCtx, k.logger))
	genesis.Params.ChainDomain = vmGenState.Params.ChainDomain
	genesis.Params.SysnamesPkgpath = vmGenState.Params.SysNamesPkgPath

	return genesis, nil
}
