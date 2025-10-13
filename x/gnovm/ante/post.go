package ante

import (
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bft "github.com/gnolang/gno/tm2/pkg/bft/types"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

type GnoTransactionsPost struct {
	logger log.Logger
	k      keeper.Keeper
}

func (gtp GnoTransactionsPost) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	if success {
		gnoCtx := gnosdk.NewContext(
			gnosdk.RunTxModeDeliver,
			nil, // MultiStore provided by keeper's VM wrapper
			&bft.Header{ChainID: ctx.ChainID()},
			types.NewSlogFromCosmosLogger(gtp.logger),
		)
		gtp.k.CommitGnoTransactionStore(gnoCtx)
	}

	return next(ctx, tx, simulate, success)
}
