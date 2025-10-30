package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

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
		msg.Package = &types.Package{
			Name: "p",
			Path: "gno.land/r/demo/p",
			Files: []*types.File{
				{
					Name: "p.gno",
					Body: "package p\n",
				},
			},
		}
		msg.Deposit = sdk.NewCoins()
		msg.MaxDeposit = sdk.NewInt64Coin("ugnot", 0)

		ms := keeper.NewMsgServerImpl(&k)
		if _, err := ms.AddPackage(ctx, msg); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "add-package executed"), nil, nil
	}
}
