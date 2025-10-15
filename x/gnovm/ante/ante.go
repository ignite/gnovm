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
			ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter()) // TODO: check if this is a hack or if the VM does prevent infinite gas.
			return ctx, nil
		}
	}

	return ctx, nil
}
