package ante

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/gnovm/x/gnovm/types"
)

func NewAnteHandler() sdk.AnteDecorator {
	return &gnoAnteHandler{}
}

type gnoAnteHandler struct{}

func (gnoAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgs := tx.GetMsgs()

	gnoVMCount := 0
	nonGnoVMCount := 0

	for _, msg := range msgs {
		switch {
		case sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&types.MsgRun{}) ||
			sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&types.MsgAddPackage{}) ||
			sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&types.MsgCall{}):
			gnoVMCount++
		default:
			nonGnoVMCount++
		}
	}

	// Reject transactions that mix GnoVM and non-GnoVM messages
	if gnoVMCount > 0 && nonGnoVMCount > 0 {
		return ctx, fmt.Errorf("cannot mix GnoVM messages with non-GnoVM messages in the same transaction")
	}

	if gnoVMCount > 0 && simulate {
		newCtx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	}

	return next(newCtx, tx, simulate)
}
