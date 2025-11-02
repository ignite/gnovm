package keeper_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/pkg/gnomod"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/std"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

// CreateMemPackageFromFiles creates a MemPackage from file contents.
func CreateMemPackageFromFiles(name, path string, files map[string]string) (*std.MemPackage, error) {
	memFiles := make([]*std.MemFile, 0, len(files))
	for filename, content := range files {
		memFiles = append(memFiles, &std.MemFile{
			Name: filename,
			Body: content,
		})
	}

	return &std.MemPackage{
		Name:  name,
		Path:  path,
		Files: memFiles,
	}, nil
}

// ReadMemPackageFromDir reads a MemPackage from a directory on disk.
func ReadMemPackageFromDir(dirPath string) (*std.MemPackage, error) {
	gnoMod, err := gnomod.ParseDir(dirPath)
	if err != nil {
		return nil, err
	}

	return gnolang.ReadMemPackage(dirPath, gnoMod.Module, gnolang.MPAnyAll)
}

// TestMsgAddPackage_Success validates adding the counter package successfully.
func TestMsgAddPackage_Success(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	creatorBytes := f.keeper.GetAuthority()
	creatorStr, err := f.addressCodec.BytesToString(creatorBytes)
	require.NoError(t, err)

	// Read the counter package from testdata directory
	testdataPath := filepath.Join("testdata", "counter")
	mpkg, err := ReadMemPackageFromDir(testdataPath)
	require.NoError(t, err)

	pkgBz, err := json.Marshal(mpkg)
	require.NoError(t, err)

	// Use sufficient deposit to cover storage costs (2949 bytes * 1 stake/byte = 2949 stake)
	deposit, _ := sdk.ParseCoinsNormalized("2949stake")
	send, _ := sdk.ParseCoinsNormalized("1000stake")
	maxDeposit, _ := sdk.ParseCoinsNormalized("5000stake")

	f.authKeeper.EXPECT().GetAccount(f.ctx, creatorBytes).
		Return(authtypes.NewBaseAccountWithAddress(creatorBytes))
	// Expected SendCoins for the send parameter
	pkgAddr := gnolang.DerivePkgCryptoAddr("gno.land/r/demo/counter").Bytes()
	f.bankKeeper.EXPECT().SendCoins(f.ctx, creatorBytes, pkgAddr, send)
	// Expected SendCoins for the storage deposit
	storageDepositAddr := gnolang.DeriveStorageDepositCryptoAddr("gno.land/r/demo/counter").Bytes()
	f.bankKeeper.EXPECT().SendCoins(f.ctx, creatorBytes, storageDepositAddr, deposit)

	msg := types.NewMsgAddPackage(creatorStr, send, maxDeposit, pkgBz)

	resp, err := ms.AddPackage(f.ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestMsgCall_Success validates calling the counter Increment function.
func TestMsgCall_Success(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	creatorBytes := f.keeper.GetAuthority()
	creatorStr, err := f.addressCodec.BytesToString(creatorBytes)
	require.NoError(t, err)

	// Read the counter package from testdata directory
	testdataPath := filepath.Join("testdata", "counter")
	mpkg, err := ReadMemPackageFromDir(testdataPath)
	require.NoError(t, err)

	pkgBz, err := json.Marshal(mpkg)
	require.NoError(t, err)

	send, _ := sdk.ParseCoinsNormalized("1000stake")
	deposit, _ := sdk.ParseCoinsNormalized("2949stake")
	maxDeposit, _ := sdk.ParseCoinsNormalized("5000stake")

	f.authKeeper.EXPECT().GetAccount(f.ctx, creatorBytes).
		Return(authtypes.NewBaseAccountWithAddress(creatorBytes)).AnyTimes()
	// Expected SendCoins for the send parameter
	pkgAddr := gnolang.DerivePkgCryptoAddr("gno.land/r/demo/counter").Bytes()
	f.bankKeeper.EXPECT().SendCoins(f.ctx, creatorBytes, pkgAddr, send)
	// Expected SendCoins for the storage deposit
	storageDepositAddr := gnolang.DeriveStorageDepositCryptoAddr("gno.land/r/demo/counter").Bytes()
	f.bankKeeper.EXPECT().SendCoins(f.ctx, creatorBytes, storageDepositAddr, deposit)

	addPkgMsg := types.NewMsgAddPackage(creatorStr, send, maxDeposit, pkgBz)
	_, err = ms.AddPackage(f.ctx, addPkgMsg)
	require.NoError(t, err)

	send, _ = sdk.ParseCoinsNormalized("2000stake")
	deposit, _ = sdk.ParseCoinsNormalized("5stake")
	amount := sdk.NewInt64Coin("ugnot", 0)

	// Expected SendCoins for the send parameter
	f.bankKeeper.EXPECT().SendCoins(f.ctx, creatorBytes, pkgAddr, send)
	// Expected SendCoins for the storage deposit
	f.bankKeeper.EXPECT().SendCoins(f.ctx, creatorBytes, storageDepositAddr, deposit)

	callMsg := types.NewMsgCall(creatorStr, send, amount, mpkg.Path, "Increment", []string{})
	resp, err := ms.Call(f.ctx, callMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Contains(t, resp.Result, "1")
}

// TestMsgRun_Success validates running a simple script.
func TestMsgRun_Success(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	callerStr, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	callerBytes, err := f.addressCodec.StringToBytes(callerStr)
	require.NoError(t, err)

	var caddr crypto.Address
	copy(caddr[:], callerBytes)
	runPath := "gno.land/e/" + caddr.String() + "/run"

	mpkg := &std.MemPackage{
		Name: "main",
		Path: runPath,
		Files: []*std.MemFile{
			{
				Name: "main.gno",
				Body: `package main

func main() {
	println("Hello, GnoVM!")
}
`,
			},
		},
	}

	pkgBz, err := json.Marshal(mpkg)
	require.NoError(t, err)

	amount := sdk.NewInt64Coin("ugnot", 0)

	f.authKeeper.EXPECT().GetAccount(f.ctx, callerBytes).
		Return(authtypes.NewBaseAccountWithAddress(callerBytes))
	f.bankKeeper.EXPECT().SendCoins(f.ctx, callerBytes, callerBytes, sdk.NewCoins())

	msg := types.NewMsgRun(callerStr, sdk.NewCoins(), amount, pkgBz)

	resp, err := ms.Run(f.ctx, msg)
	require.NoError(t, err)
	require.Equal(t, resp, &types.MsgRunResponse{
		Result: "Hello, GnoVM!\n",
	})
}

// TestMsgRun_Failed ensures MsgRun fails with a minimal invalid package.
func TestMsgRun_Failed(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	callerStr, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	callerBytes, err := f.addressCodec.StringToBytes(callerStr)
	require.NoError(t, err)

	var caddr crypto.Address
	copy(caddr[:], callerBytes)
	runPath := "gno.land/e/" + caddr.String() + "/run"

	mpkg := std.MemPackage{
		Name: "main",
		Path: runPath,
		Files: []*std.MemFile{
			{Name: "main.gno", Body: "package main\n"},
		},
	}
	pkgBz, err := json.Marshal(&mpkg)
	require.NoError(t, err)

	f.authKeeper.EXPECT().GetAccount(f.ctx, callerBytes).
		Return(authtypes.NewBaseAccountWithAddress(callerBytes))
	f.bankKeeper.EXPECT().SendCoins(f.ctx, callerBytes, callerBytes, sdk.NewCoins())

	msg := &types.MsgRun{
		Caller:     callerStr,
		Send:       sdk.NewCoins(),
		MaxDeposit: sdk.NewInt64Coin("ugnot", 0),
		Pkg:        pkgBz,
	}

	_, err = ms.Run(f.ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to run VM")
}

// TestMsgAddPackage_Failed ensures MsgAddPackage fails with a minimal invalid package.
func TestMsgAddPackage_Failed(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	creatorBytes := f.keeper.GetAuthority()
	creatorStr, err := f.addressCodec.BytesToString(creatorBytes)
	require.NoError(t, err)

	mpkg := std.MemPackage{
		Name: "p",
		Path: "gno.land/r/demo/p",
		Files: []*std.MemFile{
			{Name: "p.gno", Body: "package p\n"},
		},
	}
	pkgBz, err := json.Marshal(&mpkg)
	require.NoError(t, err)

	f.authKeeper.EXPECT().GetAccount(f.ctx, creatorBytes).
		Return(authtypes.NewBaseAccountWithAddress(creatorBytes))

	msg := &types.MsgAddPackage{
		Creator:    creatorStr,
		Send:       sdk.NewCoins(),
		MaxDeposit: sdk.NewCoins(sdk.NewInt64Coin("ugnot", 0)),
		Package:    pkgBz,
	}

	_, err = ms.AddPackage(f.ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to add package")
}

// TestMsgCall_Failed validates error handling when calling missing realm function.
func TestMsgCall_Failed(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(&f.keeper)

	require.NoError(t, f.keeper.InitGenesis(f.ctx, types.GenesisState{Params: types.DefaultParams()}))

	callerStr, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	msg := &types.MsgCall{
		Caller:     callerStr,
		Send:       sdk.NewCoins(),
		MaxDeposit: sdk.NewInt64Coin("ugnot", 0),
		PkgPath:    "gno.land/r/demo/nonexistent",
		Function:   "main",
		Args:       nil,
	}

	_, err = ms.Call(f.ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "panic while calling VM")
}
