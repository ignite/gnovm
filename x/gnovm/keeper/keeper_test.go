package keeper_test

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
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
}

// mockAuthKeeper implements types.AuthKeeper for tests.
type mockAuthKeeper struct {
	accounts map[string]sdk.AccountI
}

func newMockAuthKeeper() *mockAuthKeeper {
	return &mockAuthKeeper{accounts: make(map[string]sdk.AccountI)}
}

func (m *mockAuthKeeper) GetAccount(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	key := addr.String()
	if acc, ok := m.accounts[key]; ok {
		return acc
	}
	acc := authtypes.NewBaseAccountWithAddress(addr)
	m.accounts[key] = acc
	return acc
}

// mockBankKeeper implements types.BankKeeper for tests.
type mockBankKeeper struct {
	balances map[string]sdk.Coins
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{balances: make(map[string]sdk.Coins)}
}

func (m *mockBankKeeper) GetAllBalances(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	if c, ok := m.balances[addr.String()]; ok {
		return c
	}
	return sdk.NewCoins()
}

func (m *mockBankKeeper) SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.GetAllBalances(ctx, addr)
}

func (m *mockBankKeeper) SendCoins(_ context.Context, from, to sdk.AccAddress, amt sdk.Coins) error {
	if amt.IsAnyNegative() {
		return fmt.Errorf("negative amount")
	}
	if amt.IsZero() {
		return nil
	}
	fromKey := from.String()
	toKey := to.String()

	fromBal := m.balances[fromKey]
	if fromBal == nil {
		fromBal = sdk.NewCoins()
	}
	toBal := m.balances[toKey]
	if toBal == nil {
		toBal = sdk.NewCoins()
	}

	if !fromBal.IsAllGTE(amt) {
		return fmt.Errorf("insufficient funds")
	}
	fromBal = fromBal.Sub(amt...)
	toBal = toBal.Add(amt...)

	m.balances[fromKey] = fromBal
	m.balances[toKey] = toBal
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	panic("not implemented")
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

	k := keeper.NewKeeper(
		log.NewTestLogger(t),
		storeKey,
		memStoreKey,
		encCfg.Codec,
		addressCodec,
		authority,
		newMockAuthKeeper(),
		newMockBankKeeper(),
	)

	// Initialize params
	if err := k.Params.Set(sdkCtx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          sdkCtx,
		keeper:       k,
		addressCodec: addressCodec,
	}
}
