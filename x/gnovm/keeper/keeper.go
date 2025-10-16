package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"

	"github.com/ignite/gnovm/x/gnovm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	bft "github.com/gnolang/gno/tm2/pkg/bft/types"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	gnostore "github.com/gnolang/gno/tm2/pkg/store"
)

type Keeper struct {
	*vm.VMKeeper
	// tracks if VmKeeper has been initialized
	vmInitialized bool

	logger          log.Logger
	storeService    corestore.KVStoreService
	storeKey        *storetypes.KVStoreKey
	memStoreService corestore.MemoryStoreService
	memStoreKey     *storetypes.MemoryStoreKey
	cdc             codec.Codec
	addressCodec    address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]
	// vmKeeperParams manages VM module parameters and state.
	vmParams *vmKeeperParams

	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper
}

// NewKeeper creates a new Keeper instance.
func NewKeeper(
	logger log.Logger,
	storeKey *storetypes.KVStoreKey,
	memStoreKey *storetypes.MemoryStoreKey,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	storeService := runtime.NewKVStoreService(storeKey)
	memStoreService := runtime.NewMemStoreService(memStoreKey)

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		logger:          logger,
		storeService:    storeService,
		storeKey:        storeKey,
		memStoreService: memStoreService,
		memStoreKey:     memStoreKey,
		cdc:             cdc,
		addressCodec:    addressCodec,
		authority:       authority,
		Params:          collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		authKeeper:      authKeeper,
		bankKeeper:      bankKeeper,
	}
	k.vmParams = &vmKeeperParams{k: &k}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	// VMKeeper will be created lazily when needed
	// k.VMKeeper = nil initially

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// initializeVMKeeper creates and initializes the VMKeeper with a proper MultiStore.
// This should be called when we have access to a proper SDK context.
func (k *Keeper) initializeVMKeeper(sdkCtx sdk.Context) error {
	k.vmParams.SetSDKContext(sdkCtx)

	// check if already initialized to avoid double initialization
	if k.vmInitialized {
		return nil
	}

	// Create a MultiStore wrapper for initialization
	multiStore := NewGnovmMultiStore(
		k.logger,
		k.storeService,
		k.memStoreService,
		gnostore.NewStoreKey(k.storeKey.Name()),
		gnostore.NewStoreKey(k.memStoreKey.Name()),
	)

	chainID := sdkCtx.ChainID()
	if chainID == "" {
		chainID = "default_chain_id" // TODO: investigate why chainID is emtpy here
	}

	// Create a clean gno context for initialization
	gnoCtx := gnosdk.NewContext(
		gnosdk.RunTxModeDeliver,
		multiStore,
		&bft.Header{ChainID: chainID},
		types.NewSlogFromCosmosLogger(k.logger),
	)

	// Set the context on the multistore after creating the context
	if lazyStore, ok := multiStore.(*gnovmMultiStore); ok {
		lazyStore.SetContext(gnoCtx, sdkCtx)
	}

	// Create VMKeeper with proper initialization
	k.VMKeeper = vm.NewVMKeeper(
		k.storeKey,
		k.memStoreKey,
		vmAuthKeeper{k.logger, k.authKeeper, k.bankKeeper, k.vmParams},
		vmBankKeeper{k.logger, k.bankKeeper, k.vmParams},
		k.vmParams,
	)

	// Initialize the VMKeeper with the multistore
	k.VMKeeper.Initialize(types.NewSlogFromCosmosLogger(k.logger), multiStore)
	k.vmInitialized = true

	return nil
}

// BuildGnoContext initializes the VM (if needed), creates a Gno context using
// the MultiStore wrapper bound to the provided sdkCtx, and returns a per-tx context
// with a transaction store attached.
func (k *Keeper) BuildGnoContext(sdkCtx sdk.Context) (gnosdk.Context, error) {
	if err := k.initializeVMKeeper(sdkCtx); err != nil {
		return gnosdk.Context{}, err
	}

	var mode gnosdk.RunTxMode
	switch sdkCtx.ExecMode() {
	case sdk.ExecModeCheck, sdk.ExecModeReCheck:
		mode = gnosdk.RunTxModeCheck
	case sdk.ExecModeSimulate:
		mode = gnosdk.RunTxModeSimulate
	case sdk.ExecModeFinalize:
		mode = gnosdk.RunTxModeDeliver
	default:
		return gnosdk.Context{}, fmt.Errorf("invalid exec mode: %v", sdkCtx.ExecMode())
	}

	// Create MultiStore wrapper for transaction
	ms := NewGnovmMultiStore(
		k.logger,
		k.storeService,
		k.memStoreService,
		gnostore.NewStoreKey(k.storeKey.Name()),
		gnostore.NewStoreKey(k.memStoreKey.Name()),
	)

	chainID := sdkCtx.ChainID()
	if chainID == "" {
		chainID = "default_chain_id" // TODO: investigate why chainID is emtpy here
	}

	gnoCtx := gnosdk.NewContext(
		mode,
		ms,
		&bft.Header{ChainID: chainID},
		types.NewSlogFromCosmosLogger(k.logger),
	)

	// Set the context on the multistore after creating the context
	if lazyStore, ok := ms.(*gnovmMultiStore); ok {
		lazyStore.SetContext(gnoCtx, sdkCtx)
	}

	gnoCtx = k.VMKeeper.MakeGnoTransactionStore(gnoCtx)
	return gnoCtx, nil
}
