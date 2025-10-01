package ante

import (
	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

type GnoTransactionsAnte struct {
	logger log.Logger
	k      keeper.Keeper
}

func (gta GnoTransactionsAnte) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	_ = gta.k.VMKeeper.MakeGnoTransactionStore(types.GnoContextFromSDKContext(ctx, gta.logger))

	return next(ctx, tx, simulate)
}
