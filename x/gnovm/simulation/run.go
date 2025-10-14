package simulation

import (
	"encoding/hex"
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

func SimulateMsgRun(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRun{
			Caller: simAccount.Address.String(),
		}

		// Build a minimal valid in-memory package for MsgRun
		files := []*std.MemFile{
			{
				Name: "main.gno",
				Body: "package main\n",
			},
		}

		// Expected run path is gno.land/e/<caller_hex>/run where caller is the crypto.Address string form.
		// We approximate it using the hex form of the account bytes to satisfy MemPackage path validation.
		addrHex := hex.EncodeToString(simAccount.Address.Bytes())

		mpkg := std.MemPackage{
			Name:  "main",
			Path:  "gno.land/e/" + addrHex + "/run",
			Files: files,
		}

		bz, _ := json.Marshal(&mpkg)
		msg.Pkg = bz

		// Execute through the message server
		ms := keeper.NewMsgServerImpl(&k)
		_, err := ms.Run(ctx, msg)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "run executed"), nil, nil
	}
}
