package ante

import (
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/gnovm/x/gnovm/keeper"
)

func NewTransactionsPostHandler(logger log.Logger, k keeper.Keeper) sdk.PostDecorator {
	return &gnoTransactionsPost{
		logger: logger,
		k:      k,
	}
}

type gnoTransactionsPost struct {
	logger log.Logger
	k      keeper.Keeper
}

func (gtp *gnoTransactionsPost) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	if success {
		gnoCtx, err := gtp.k.BuildGnoContextWithStore(ctx)
		if err != nil {
			return ctx, err
		}
		gtp.k.CommitGnoTransactionStore(gnoCtx)
	}

	return next(ctx, tx, simulate, success)
}
