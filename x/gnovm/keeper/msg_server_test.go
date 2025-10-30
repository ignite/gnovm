package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gnolang/gno/tm2/pkg/crypto"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

// TestMsgRun_Failed ensures MsgRun fails with a minimal valid in-memory package.
func TestMsgRun_Failed(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	// Initialize VM genesis params (chain domain, etc.) before executing messages
	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	// Use module authority as a valid caller address.
	callerStr, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	// Build the expected run path: gno.land/e/<caller-crypto-addr>/run
	callerBytes, err := f.addressCodec.StringToBytes(callerStr)
	require.NoError(t, err)
	var caddr crypto.Address
	copy(caddr[:], callerBytes)
	runPath := "gno.land/e/" + caddr.String() + "/run"

	// Minimal valid MemPackage
	pkg := &types.Package{
		Name: "main",
		Path: runPath,
		Files: []*types.File{
			{Name: "main.gno", Body: "package main\n"},
		},
	}
	// setup mock expectations
	f.authKeeper.EXPECT().GetAccount(f.ctx, callerBytes).
		Return(authtypes.NewBaseAccountWithAddress(callerBytes))
	// MsgRun transfers send from caller to caller (no-op)
	f.bankKeeper.EXPECT().SendCoins(f.ctx, callerBytes, callerBytes,
		sdk.NewCoins())

	msg := &types.MsgRun{
		Caller:     callerStr,
		Send:       sdk.NewCoins(),
		MaxDeposit: sdk.NewInt64Coin("ugnot", 0),
		Pkg:        pkg,
	}

	_, err = ms.Run(f.ctx, msg)
	require.Error(t, err)
	println(err.Error())
	require.Contains(t, err.Error(), "failed to run VM")
}

// TestMsgAddPackage_Failed ensures MsgAddPackage fails with a minimal valid package.
func TestMsgAddPackage_Failed(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	// Initialize VM genesis params before executing messages
	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	creatorBytes := f.keeper.GetAuthority()
	creatorStr, err := f.addressCodec.BytesToString(creatorBytes)
	require.NoError(t, err)

	// Minimal valid package for add-package
	pkg := &types.Package{
		Name: "p",
		Path: "gno.land/r/demo/p",
		Files: []*types.File{
			{Name: "p.gno", Body: "package p\n"},
		},
	}
	// setup mock expectations
	f.authKeeper.EXPECT().GetAccount(f.ctx, creatorBytes).
		Return(authtypes.NewBaseAccountWithAddress(creatorBytes))

	msg := &types.MsgAddPackage{
		Creator:    creatorStr,
		Deposit:    sdk.NewCoins(),
		MaxDeposit: sdk.NewInt64Coin("ugnot", 0),
		Package:    pkg,
	}

	_, err = ms.AddPackage(f.ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to add package")
}

// TestMsgCall_Failed validates forwarding to VMKeeper and error wrapping on missing realm/function.
func TestMsgCall_Failed(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	// Initialize VM genesis params before executing messages
	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	callerStr, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	// Provide valid fields; underlying VM likely errors due to missing realm/function.
	msg := &types.MsgCall{
		Caller:     callerStr,
		Send:       sdk.NewCoins(),
		MaxDeposit: sdk.NewInt64Coin("ugnot", 0),
		PkgPath:    "gno.land/r/demo/p",
		Function:   "main",
		Args:       nil,
	}

	_, err = ms.Call(f.ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "panic while calling VM")
}
