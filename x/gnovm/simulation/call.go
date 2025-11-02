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

func SimulateMsgCall(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCall{
			Caller: simAccount.Address.String(),
		}

		// Fill required fields to pass validation and execute via msg server
		msg.PkgPath = "gno.land/r/demo/p"
		msg.Function = "main"
		msg.Args = nil
		msg.MaxDeposit = sdk.NewCoins(sdk.NewInt64Coin("test", 0))

		ms := keeper.NewMsgServerImpl(&k)
		if _, err := ms.Call(ctx, msg); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "call executed"), nil, nil
	}
}
