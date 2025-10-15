package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/sdk/params"
	"github.com/ignite/gnovm/x/gnovm/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Ensure VMKeeper is initialized before using it
	if err := k.initializeVMKeeper(sdkCtx); err != nil {
		return err
	}

	// Create a safe gno context for genesis operations
	gnoCtx, err := k.BuildGnoContextWithStore(sdkCtx)
	if err != nil {
		return err
	}

	realmParams := make([]params.Param, 0)
	if len(genState.RealmParams) > 0 {
		if err := json.Unmarshal(genState.RealmParams, &realmParams); err != nil {
			return err
		}
	}

	// Initialize the VMKeeper with the genesis state
	k.VMKeeper.InitGenesis(
		gnoCtx,
		vm.GenesisState{
			Params:      genState.Params.ToVmParams(),
			RealmParams: realmParams,
		},
	)

	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	// Ensure VMKeeper is initialized before using it
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.initializeVMKeeper(sdkCtx); err != nil {
		return nil, err
	}

	// Create a safe gno context for genesis operations
	gnoCtx, err := k.BuildGnoContextWithStore(sdkCtx)
	if err != nil {
		return nil, err
	}

	// no need check to module params state, as it is in sync with the VMKeeper
	vmGenState := k.VMKeeper.ExportGenesis(gnoCtx)

	realmParams, err := json.Marshal(vmGenState.RealmParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal realm params: %w", err)
	}

	genesis := types.DefaultGenesis()
	genesis.Params = types.VmParamsToParams(vmGenState.Params)
	genesis.RealmParams = realmParams

	return genesis, nil
}
