package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

type GnoTransactionsPost struct {
	k keeper.Keeper
}

func (gtp GnoTransactionsPost) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	if success {
		gtp.k.CommitGnoTransactionStore(types.GnoContextFromSDKContext(ctx))
	}

	return next(ctx, tx, simulate, success)
}
