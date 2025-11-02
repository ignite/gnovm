package simulation

import (
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func SimulateMsgAddPackage(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAddPackage{
			Creator: simAccount.Address.String(),
		}

		// Build a minimal valid package and execute via MsgServer
		files := []*std.MemFile{
			{
				Name: "p.gno",
				Body: "package p\n",
			},
		}
		mpkg := std.MemPackage{
			Name:  "p",
			Path:  "gno.land/r/demo/p",
			Files: files,
		}
		bz, _ := json.Marshal(&mpkg)
		msg.Package = bz
		msg.Send = sdk.NewCoins()
		msg.MaxDeposit = sdk.NewCoins(sdk.NewInt64Coin("ugnot", 0))

		ms := keeper.NewMsgServerImpl(&k)
		if _, err := ms.AddPackage(ctx, msg); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "add-package executed"), nil, nil
	}
}
