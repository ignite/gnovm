package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/ignite/gnovm/x/gnovm/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	chainID := sdkCtx.ChainID()
	if chainID == "" {
		return errors.New("chainID is empty")
	}

	// Ensure VMKeeper is initialized before using it
	if err := k.initializeVMKeeper(sdkCtx); err != nil {
		return err
	}

	// Create a safe gno context for genesis operations
	gnoCtx, err := k.BuildGnoContextWithStore(sdkCtx)
	if err != nil {
		return err
	}
	defer k.VMKeeper.CommitGnoTransactionStore(gnoCtx)

	// todo: module params from the module itself and from the vmkeeper must stay in sync
	k.VMKeeper.InitGenesis(
		gnoCtx,
		vm.GenesisState{
			Params: genState.Params.ToVmParams(),
		},
	)

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	chainID := sdkCtx.ChainID()
	if chainID == "" {
		return nil, errors.New("chainID is empty")
	}

	// Ensure VMKeeper is initialized before using it
	if err := k.initializeVMKeeper(sdkCtx); err != nil {
		return nil, err
	}

	// Create a safe gno context for genesis operations
	gnoCtx, err := k.BuildGnoContextWithStore(sdkCtx)
	if err != nil {
		return nil, err
	}
	defer k.VMKeeper.CommitGnoTransactionStore(gnoCtx)

	vmGenState := k.VMKeeper.ExportGenesis(gnoCtx)
	genesis.Params.ChainDomain = vmGenState.Params.ChainDomain
	genesis.Params.SysnamesPkgpath = vmGenState.Params.SysNamesPkgPath

	return genesis, nil
}
