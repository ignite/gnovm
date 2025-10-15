package ante

import (
	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/keeper"
)

func NewTransactionsAnteHandler(logger log.Logger, k keeper.Keeper) sdk.AnteDecorator {
	return &gnoTransactionsAnte{
		logger: logger,
		k:      k,
	}
}

type gnoTransactionsAnte struct {
	logger log.Logger
	k      keeper.Keeper
}

func (gta gnoTransactionsAnte) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	gnoCtx, err := gta.k.BuildGnoContextWithStore(ctx)
	if err != nil {
		return ctx, err
	}
	_ = gta.k.VMKeeper.MakeGnoTransactionStore(gnoCtx) // TODO

	return next(ctx, tx, simulate)
}
