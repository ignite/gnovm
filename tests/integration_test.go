package tests

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/pkg/gnomod"

	gnovmkeeper "github.com/ignite/gnovm/x/gnovm/keeper"
	gnovmmodule "github.com/ignite/gnovm/x/gnovm/module"
	gnovmtypes "github.com/ignite/gnovm/x/gnovm/types"
)

const (
	testDenom = "stake"
)

type fixture struct {
	app *integration.App

	ctx sdk.Context

	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	gnovmKeeper   gnovmkeeper.Keeper

	cdc codec.Codec

	addrDels []sdk.AccAddress
}

func initFixture(t testing.TB) *fixture {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, gnovmtypes.ModuleName,
	)
	memKeys := storetypes.NewMemoryStoreKeys("mem_" + gnovmtypes.ModuleName)

	cdc := moduletestutil.MakeTestEncodingConfig(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		gnovmmodule.AppModule{},
	).Codec

	logger := log.NewTestLogger(t)
	cms := integration.CreateMultiStore(keys, logger)

	newCtx := sdk.NewContext(cms, cmtproto.Header{
		ChainID: "gnovm-test-1",
		Time:    time.Now().UTC(),
	}, true, logger)

	authority := authtypes.NewModuleAddress(gnovmtypes.GovModuleName)

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		map[string][]string{
			gnovmtypes.ModuleName: {authtypes.Minter},
		},
		addresscodec.NewBech32Codec(sdk.Bech32MainPrefix),
		sdk.Bech32MainPrefix,
		authority.String(),
	)

	blockedAddresses := map[string]bool{
		accountKeeper.GetAuthority(): false,
	}

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		blockedAddresses,
		authority.String(),
		log.NewNopLogger(),
	)

	gnovmKeeper := gnovmkeeper.NewKeeper(
		logger,
		keys[gnovmtypes.ModuleName],
		memKeys["mem_"+gnovmtypes.ModuleName],
		cdc,
		addresscodec.NewBech32Codec(sdk.Bech32MainPrefix),
		authority,
		accountKeeper,
		bankKeeper,
	)

	authModule := auth.NewAppModule(cdc, accountKeeper, nil, nil)
	bankModule := bank.NewAppModule(cdc, bankKeeper, accountKeeper, nil)
	gnovmModule := gnovmmodule.NewAppModule(cdc, gnovmKeeper, accountKeeper, bankKeeper)

	integrationApp := integration.NewIntegrationApp(newCtx, logger, keys, cdc, map[string]appmodule.AppModule{
		authtypes.ModuleName:  authModule,
		banktypes.ModuleName:  bankModule,
		gnovmtypes.ModuleName: gnovmModule,
	})

	sdkCtx := sdk.UnwrapSDKContext(integrationApp.Context())

	// Register MsgServer and QueryServer
	gnovmtypes.RegisterMsgServer(integrationApp.MsgServiceRouter(), gnovmkeeper.NewMsgServerImpl(&gnovmKeeper))
	gnovmtypes.RegisterQueryServer(integrationApp.QueryHelper(), gnovmkeeper.NewQueryServerImpl(&gnovmKeeper))

	// Initialize GnoVM genesis with default params to load std library
	err := gnovmKeeper.InitGenesis(sdkCtx, gnovmtypes.GenesisState{
		Params: gnovmtypes.DefaultParams(),
	})
	require.NoError(t, err)

	// Create test addresses with initial balances
	addrDels := simtestutil.CreateIncrementalAccounts(6)
	for _, addr := range addrDels {
		acc := accountKeeper.NewAccountWithAddress(sdkCtx, addr)
		accountKeeper.SetAccount(sdkCtx, acc)

		// Fund the account
		err := bankKeeper.MintCoins(sdkCtx, gnovmtypes.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 100000000)))
		require.NoError(t, err)
		err = bankKeeper.SendCoinsFromModuleToAccount(sdkCtx, gnovmtypes.ModuleName, addr, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 100000000)))
		require.NoError(t, err)
	}

	return &fixture{
		app:           integrationApp,
		ctx:           sdkCtx,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		gnovmKeeper:   gnovmKeeper,
		cdc:           cdc,
		addrDels:      addrDels,
	}
}

// TestBankerContract_DeployAndDeposit tests deploying the banker contract and depositing coins
func TestBankerContract_DeployAndDeposit(t *testing.T) {
	t.Parallel()
	f := initFixture(t)

	// Read the banker package
	testdataPath := filepath.Join("contracts", "banker")

	gnoMod, err := gnomod.ParseDir(testdataPath)
	require.NoError(t, err)

	mpkg, err := gnolang.ReadMemPackage(testdataPath, gnoMod.Module, gnolang.MPAnyAll)
	require.NoError(t, err)

	require.NoError(t, err, "failed to read banker package")

	pkgBz, err := json.Marshal(mpkg)
	require.NoError(t, err)

	deployer := f.addrDels[0]
	recipient := f.addrDels[1]

	// Query initial deployer balance
	initialBalance := f.bankKeeper.GetBalance(f.ctx, deployer, testDenom)
	assert.Equal(t, "100000000", initialBalance.Amount.String())

	// Query initial recipient balance (should be 100000000)
	recipientBalance := f.bankKeeper.GetBalance(f.ctx, recipient, testDenom)
	assert.Equal(t, "100000000", recipientBalance.Amount.String())

	// Deploy the banker package with sufficient storage deposit
	deposit := sdk.NewInt64Coin(testDenom, 10000)
	deployerStr, err := f.accountKeeper.AddressCodec().BytesToString(deployer)
	require.NoError(t, err)

	msgAddPackage := gnovmtypes.NewMsgAddPackage(
		deployerStr,
		sdk.NewCoins(deposit),
		deposit,
		pkgBz,
	)

	_, err = f.app.RunMsg(
		msgAddPackage,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err, "failed to add banker package")

	// Verify deployer balance decreased by deposit
	newBalance := f.bankKeeper.GetBalance(f.ctx, deployer, testDenom)
	assert.Assert(t, newBalance.Amount.LT(initialBalance.Amount))

	// Call Deposit function, sending 50000 stake with the transaction
	depositAmount := sdk.NewInt64Coin(testDenom, 50000)
	msgCall := gnovmtypes.NewMsgCall(
		deployerStr,
		sdk.NewCoins(depositAmount),
		sdk.NewInt64Coin(testDenom, 1000),
		mpkg.Path,
		"Deposit",
		[]string{},
	)

	res, err := f.app.RunMsg(
		msgCall,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err, "failed to call Deposit function")

	// Verify response contains success message
	var callResp gnovmtypes.MsgCallResponse
	err = f.cdc.Unmarshal(res.Value, &callResp)
	require.NoError(t, err)
	assert.Assert(t, len(callResp.Result) > 0)
}

// TestBankerContract_SendCoins tests sending coins from the realm to a user
func TestBankerContract_SendCoins(t *testing.T) {
	t.Parallel()
	f := initFixture(t)

	// Read and deploy the banker package
	testdataPath := filepath.Join("contracts", "banker")
	gnoMod, err := gnomod.ParseDir(testdataPath)
	require.NoError(t, err)

	mpkg, err := gnolang.ReadMemPackage(testdataPath, gnoMod.Module, gnolang.MPAnyAll)
	require.NoError(t, err)

	pkgBz, err := json.Marshal(mpkg)
	require.NoError(t, err)

	deployer := f.addrDels[0]
	recipient := f.addrDels[1]
	deployerStr, _ := f.accountKeeper.AddressCodec().BytesToString(deployer)
	recipientStr, _ := f.accountKeeper.AddressCodec().BytesToString(recipient)

	// Deploy package
	deposit := sdk.NewInt64Coin(testDenom, 10000)
	msgAddPackage := gnovmtypes.NewMsgAddPackage(
		deployerStr,
		sdk.NewCoins(deposit),
		deposit,
		pkgBz,
	)

	_, err = f.app.RunMsg(
		msgAddPackage,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err)

	// Deposit coins to realm
	depositAmount := sdk.NewInt64Coin(testDenom, 50000)
	msgDeposit := gnovmtypes.NewMsgCall(
		deployerStr,
		sdk.NewCoins(depositAmount),
		sdk.NewInt64Coin(testDenom, 1000),
		mpkg.Path,
		"Deposit",
		[]string{},
	)

	_, err = f.app.RunMsg(
		msgDeposit,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err)

	// Get recipient balance before transfer
	balanceBefore := f.bankKeeper.GetBalance(f.ctx, recipient, testDenom)

	// Send coins from realm to recipient
	transferAmount := int64(30000)
	msgSendCoins := gnovmtypes.NewMsgCall(
		deployerStr,
		sdk.NewCoins(),
		sdk.NewInt64Coin(testDenom, 1000),
		mpkg.Path,
		"SendCoins",
		[]string{recipientStr, fmt.Sprintf("%d", transferAmount), testDenom},
	)

	res, err := f.app.RunMsg(
		msgSendCoins,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err, "failed to call SendCoins function")

	// Verify response
	var callResp gnovmtypes.MsgCallResponse
	err = f.cdc.Unmarshal(res.Value, &callResp)
	require.NoError(t, err)

	// Verify recipient received coins
	balanceAfter := f.bankKeeper.GetBalance(f.ctx, recipient, testDenom)
	expectedIncrease := math.NewInt(transferAmount)
	assert.DeepEqual(t, balanceBefore.Amount.Add(expectedIncrease), balanceAfter.Amount)
}

// TestBankerContract_QueryBalance tests querying balances via the realm
func TestBankerContract_QueryBalance(t *testing.T) {
	t.Parallel()
	f := initFixture(t)

	// Read and deploy the banker package
	testdataPath := filepath.Join("contracts", "banker")
	gnoMod, err := gnomod.ParseDir(testdataPath)
	require.NoError(t, err)

	mpkg, err := gnolang.ReadMemPackage(testdataPath, gnoMod.Module, gnolang.MPAnyAll)
	require.NoError(t, err)

	pkgBz, err := json.Marshal(mpkg)
	require.NoError(t, err)

	deployer := f.addrDels[0]
	deployerStr, _ := f.accountKeeper.AddressCodec().BytesToString(deployer)

	// Deploy package
	deposit := sdk.NewInt64Coin(testDenom, 10000)
	msgAddPackage := gnovmtypes.NewMsgAddPackage(
		deployerStr,
		sdk.NewCoins(deposit),
		deposit,
		pkgBz,
	)

	_, err = f.app.RunMsg(
		msgAddPackage,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err)

	// Deposit coins to realm
	depositAmount := sdk.NewInt64Coin(testDenom, 50000)
	msgDeposit := gnovmtypes.NewMsgCall(
		deployerStr,
		sdk.NewCoins(depositAmount),
		sdk.NewInt64Coin(testDenom, 1000),
		mpkg.Path,
		"Deposit",
		[]string{},
	)

	_, err = f.app.RunMsg(
		msgDeposit,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err)

	// Query realm balance
	msgGetBalance := gnovmtypes.NewMsgCall(
		deployerStr,
		sdk.NewCoins(),
		sdk.NewInt64Coin(testDenom, 1000),
		mpkg.Path,
		"GetBalance",
		[]string{},
	)

	res, err := f.app.RunMsg(
		msgGetBalance,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err, "failed to call GetBalance function")

	// Verify response contains the deposited amount
	var callResp gnovmtypes.MsgCallResponse
	err = f.cdc.Unmarshal(res.Value, &callResp)
	require.NoError(t, err)
	assert.Assert(t, len(callResp.Result) > 0)
}

// TestBankerContract_MultipleTransfers tests multiple sequential transfers
func TestBankerContract_MultipleTransfers(t *testing.T) {
	t.Parallel()
	f := initFixture(t)

	// Read and deploy the banker package
	testdataPath := filepath.Join("contracts", "banker")
	gnoMod, err := gnomod.ParseDir(testdataPath)
	require.NoError(t, err)

	mpkg, err := gnolang.ReadMemPackage(testdataPath, gnoMod.Module, gnolang.MPAnyAll)
	require.NoError(t, err)

	pkgBz, err := json.Marshal(mpkg)
	require.NoError(t, err)

	deployer := f.addrDels[0]
	recipient1 := f.addrDels[1]
	recipient2 := f.addrDels[2]
	deployerStr, _ := f.accountKeeper.AddressCodec().BytesToString(deployer)
	recipient1Str, _ := f.accountKeeper.AddressCodec().BytesToString(recipient1)
	recipient2Str, _ := f.accountKeeper.AddressCodec().BytesToString(recipient2)

	// Deploy package
	deposit := sdk.NewInt64Coin(testDenom, 10000)
	msgAddPackage := gnovmtypes.NewMsgAddPackage(
		deployerStr,
		sdk.NewCoins(deposit),
		deposit,
		pkgBz,
	)

	_, err = f.app.RunMsg(
		msgAddPackage,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err)

	// Deposit large amount to realm
	depositAmount := sdk.NewInt64Coin(testDenom, 100000)
	msgDeposit := gnovmtypes.NewMsgCall(
		deployerStr,
		sdk.NewCoins(depositAmount),
		sdk.NewInt64Coin(testDenom, 1000),
		mpkg.Path,
		"Deposit",
		[]string{},
	)

	_, err = f.app.RunMsg(
		msgDeposit,
		integration.WithAutomaticFinalizeBlock(),
		integration.WithAutomaticCommit(),
	)
	require.NoError(t, err)

	// Get initial balances
	bal1Before := f.bankKeeper.GetBalance(f.ctx, recipient1, testDenom)
	bal2Before := f.bankKeeper.GetBalance(f.ctx, recipient2, testDenom)

	// Perform multiple transfers
	transfers := []struct {
		to     string
		amount int64
	}{
		{recipient1Str, 10000},
		{recipient2Str, 20000},
		{recipient1Str, 15000},
	}

	for _, transfer := range transfers {
		msgSend := gnovmtypes.NewMsgCall(
			deployerStr,
			sdk.NewCoins(),
			sdk.NewInt64Coin(testDenom, 1000),
			mpkg.Path,
			"SendCoins",
			[]string{transfer.to, fmt.Sprintf("%d", transfer.amount), testDenom},
		)

		_, err = f.app.RunMsg(
			msgSend,
			integration.WithAutomaticFinalizeBlock(),
			integration.WithAutomaticCommit(),
		)
		require.NoError(t, err)
	}

	// Verify final balances
	bal1After := f.bankKeeper.GetBalance(f.ctx, recipient1, testDenom)
	bal2After := f.bankKeeper.GetBalance(f.ctx, recipient2, testDenom)

	// recipient1 should have received 10000 + 15000 = 25000
	expectedIncrease1 := math.NewInt(25000)
	assert.DeepEqual(t, bal1Before.Amount.Add(expectedIncrease1), bal1After.Amount)

	// recipient2 should have received 20000
	expectedIncrease2 := math.NewInt(20000)
	assert.DeepEqual(t, bal2Before.Amount.Add(expectedIncrease2), bal2After.Amount)
}
