package ante

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func NewAnteHandler() sdk.AnteDecorator {
	return &gnoAnteHandler{}
}

type gnoAnteHandler struct{}

func (*gnoAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch {
		case sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&types.MsgRun{}) ||
			sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&types.MsgAddPackage{}) ||
			sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&types.MsgCall{}):
			// Use infinite gas meter for GnoVM transactions because the VM has its own
			// internal gas tracking mechanism through the transaction store. Using the
			// Cosmos SDK gas meter would result in double-counting gas consumption.
			// The VM prevents infinite execution through its own gas limits.
			ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
			return ctx, nil
		}
	}

	return ctx, nil
}
