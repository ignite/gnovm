package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	bft "github.com/gnolang/gno/tm2/pkg/bft/types"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/ignite/gnovm/x/gnovm/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.initializeVMKeeper(sdkCtx); err != nil {
		return err
	}

	chainID := sdkCtx.ChainID()
	if chainID == "" {
		return errors.New("chainID is empty")
	}

	// Create a safe gno context for genesis operations
	gnoCtx := gnosdk.NewContext(
		gnosdk.RunTxModeDeliver,
		nil, // multistore - VMKeeper will use our wrapper
		&bft.Header{ChainID: chainID},
		types.NewSlogFromCosmosLogger(k.logger),
	)

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
	if err := k.initializeVMKeeper(sdkCtx); err != nil {
		return nil, err
	}

	chainID := sdkCtx.ChainID()
	if chainID == "" {
		return nil, errors.New("chainID is empty")
	}

	// Create a safe gno context for genesis operations
	gnoCtx := gnosdk.NewContext(
		gnosdk.RunTxModeDeliver,
		nil, // multistore - VMKeeper will use our wrapper
		&bft.Header{ChainID: chainID},
		types.NewSlogFromCosmosLogger(k.logger),
	)

	vmGenState := k.VMKeeper.ExportGenesis(gnoCtx)
	genesis.Params.ChainDomain = vmGenState.Params.ChainDomain
	genesis.Params.SysnamesPkgpath = vmGenState.Params.SysNamesPkgPath

	return genesis, nil
}
