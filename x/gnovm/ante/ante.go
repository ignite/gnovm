package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

type GnoTransactionsAnte struct {
	k keeper.Keeper
}

func (gta GnoTransactionsAnte) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	_ = gta.k.VMKeeper.MakeGnoTransactionStore(types.GnoContextFromSDKContext(ctx))

	return next(ctx, tx, simulate)
}
