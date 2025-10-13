package ante

import (
	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bft "github.com/gnolang/gno/tm2/pkg/bft/types"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

type GnoTransactionsAnte struct {
	logger log.Logger
	k      keeper.Keeper
}

func (gta GnoTransactionsAnte) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	gnoCtx := gnosdk.NewContext(
		gnosdk.RunTxModeDeliver,
		nil, // MultiStore provided by keeper's VM wrapper
		&bft.Header{ChainID: ctx.ChainID()},
		types.NewSlogFromCosmosLogger(gta.logger),
	)
	_ = gta.k.VMKeeper.MakeGnoTransactionStore(gnoCtx) // TODO

	return next(ctx, tx, simulate)
}
