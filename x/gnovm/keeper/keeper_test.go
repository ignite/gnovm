package keeper_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	module "github.com/ignite/gnovm/x/gnovm/module"
	"github.com/ignite/gnovm/x/gnovm/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	storeService store.KVStoreService
	authKeeper   *MockAuthKeeper
	bankKeeper   *MockBankKeeper
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey("gnovm")
	memStoreKey := storetypes.NewMemoryStoreKey("memory:gnovm")

	tKey := storetypes.NewTransientStoreKey("transient_test")
	sdkCtx := testutil.DefaultContextWithDB(t, storeKey, tKey).Ctx
	sdkCtx = sdkCtx.WithChainID("gnovm-test")

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	ctrl := gomock.NewController(t)
	authKeeper := NewMockAuthKeeper(ctrl)
	bankKeeper := NewMockBankKeeper(ctrl)

	k := keeper.NewKeeper(
		log.NewTestLogger(t),
		storeKey,
		memStoreKey,
		encCfg.Codec,
		addressCodec,
		authority,
		authKeeper,
		bankKeeper,
	)

	// Initialize params
	if err := k.Params.Set(sdkCtx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          sdkCtx,
		keeper:       k,
		addressCodec: addressCodec,
		authKeeper:   authKeeper,
		bankKeeper:   bankKeeper,
		storeService: runtime.NewKVStoreService(storeKey),
	}
}
