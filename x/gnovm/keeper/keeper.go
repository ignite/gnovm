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

	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
)

type Keeper struct {
	*vm.VMKeeper

	logger       log.Logger
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper
}

func NewKeeper(
	logger log.Logger,
	storeKey *storetypes.KVStoreKey,
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
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		logger:       logger,
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		authKeeper:   authKeeper,
		bankKeeper:   bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	k.VMKeeper = vm.NewVMKeeper(
		storeKey, // TODO(@julienrbrt): possible use another one.
		storeKey,
		vmAuthKeeper{k.logger, k.authKeeper, k.bankKeeper},
		vmBankKeeper{k.logger, k.bankKeeper},
		&vmKeeperParams{&k},
	)

	k.VMKeeper.Initialize(types.NewSlogFromCosmosLogger(k.logger), nil)

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
