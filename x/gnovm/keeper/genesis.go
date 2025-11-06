package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/ignite/gnovm/x/gnovm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/sdk/params"

	"golang.org/x/tools/go/packages"
)

const defaultStdLibs = "github.com/gnolang/gno/gnovm/stdlibs"

// InitGenesis initializes the module's state from a provided genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Import all key-value pairs into the store
	if len(genState.State) > 0 {
		store := k.storeService.OpenKVStore(sdkCtx)
		for _, kv := range genState.State {
			if err := store.Set(kv.Key, kv.Value); err != nil {
				return fmt.Errorf("failed to set key-value pair during genesis import: %w", err)
			}
		}
	}

	gnoCtx, err := k.BuildGnoContext(sdkCtx)
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

	pkg, err := packages.Load(&packages.Config{Mode: packages.LoadFiles}, defaultStdLibs)
	if err != nil {
		return fmt.Errorf("failed to load gno stdlib packages: %w", err)
	}
	if len(pkg) == 0 {
		return fmt.Errorf("no gno stdlib packages found")
	}

	// Initialize the standard library
	k.VMKeeper.LoadStdlib(gnoCtx, pkg[0].Dir)
	k.VMKeeper.CommitGnoTransactionStore(gnoCtx)

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	gnoCtx, err := k.BuildGnoContext(sdkCtx)
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

	// Export all key-value pairs from the store
	store := k.storeService.OpenKVStore(sdkCtx)
	iterator, err := store.Iterator(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create store iterator: %w", err)
	}
	defer iterator.Close()

	var state []types.KVPair
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()

		// Store the key-value pair (make copies to avoid aliasing)
		state = append(state, types.KVPair{
			Key:   bytes.Clone(key),
			Value: bytes.Clone(value),
		})
	}

	genesis.State = state

	return genesis, nil
}
