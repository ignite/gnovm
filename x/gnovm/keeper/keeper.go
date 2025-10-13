package keeper

import (
	"errors"
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

	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper

	// Reference to vmKeeperParams for setting SDK context
	vmParams *vmKeeperParams
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

	// gno keeper
	k.VMKeeper = vm.NewVMKeeper(
		storeKey,
		memStoreKey,
		vmAuthKeeper{k.logger, k.authKeeper, k.bankKeeper, k.vmParams},
		vmBankKeeper{k.logger, k.bankKeeper, k.vmParams},
		k.vmParams,
	)

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// initializeVMKeeper initializes the VMKeeper with a proper MultiStore.
// This should be called when we have access to a proper SDK context.
func (k *Keeper) initializeVMKeeper(sdkCtx sdk.Context) error {
	if k.VMKeeper == nil {
		return errors.New("VMKeeper not created")
	}
	k.vmParams.SetSDKContext(sdkCtx)

	// check if already initialized to avoid double initialization
	if k.vmInitialized {
		return nil
	}

	// Create a safe gno context for the MultiStore wrapper
	gnoCtx := gnosdk.NewContext(
		gnosdk.RunTxModeDeliver,
		nil, // multistore - we'll provide our own wrapper
		&bft.Header{ChainID: "gnovm-chain"},
		types.NewSlogFromCosmosLogger(k.logger),
	)

	// Create a MultiStore wrapper that restricts access to only the gnovm store
	multiStore := NewGnovmMultiStore(
		k.logger,
		k.storeService,
		k.memStoreService,
		gnostore.NewStoreKey(k.storeKey.Name()),
		gnostore.NewStoreKey(k.memStoreKey.Name()),
		gnoCtx,
		sdkCtx,
	)

	k.VMKeeper.Initialize(types.NewSlogFromCosmosLogger(k.logger), multiStore)
	k.vmInitialized = true

	return nil
}
