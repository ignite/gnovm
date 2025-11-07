package ante

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gnostore "github.com/gnolang/gno/tm2/pkg/store"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func NewAnteHandler(k *keeper.Keeper) sdk.AnteDecorator {
	return &gnoAnteHandler{keeper: k}
}

type gnoAnteHandler struct {
	keeper *keeper.Keeper
}

func (gad *gnoAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
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
	// This is necessary because they use different gas accounting mechanisms
	if gnoVMCount > 0 && nonGnoVMCount > 0 {
		return ctx, fmt.Errorf("cannot mix GnoVM messages with non-GnoVM messages in the same transaction")
	}

	if gnoVMCount == 0 {
		return next(ctx, tx, simulate)
	}

	// Get the gas limit from the current SDK context
	// This was set by earlier ante handlers based on the transaction's gas wanted
	gasLimit := ctx.GasMeter().Limit()

	// Use infinite gas meter in the SDK context to prevent double-counting
	// The GnoVM has its own internal gas tracking mechanism through the transaction store
	newCtx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	gnoCtx, err := gad.keeper.BuildGnoContext(newCtx)
	if err != nil {
		return ctx, fmt.Errorf("failed to build gno context: %w", err)
	}

	if !simulate {
		gnoCtx = gnoCtx.WithGasMeter(gnostore.NewGasMeter(gnostore.Gas(gasLimit)))
	}

	// Set up defer/recover to handle gas consumption and out-of-gas errors
	// This pattern follows the Gno tm2/pkg/sdk/auth/ante.go implementation
	defer func() {
		if r := recover(); r != nil {
			switch ex := r.(type) {
			case gnostore.OutOfGasError:
				gasConsumed := gnoCtx.GasMeter().GasConsumed()

				log := fmt.Sprintf(
					"out of gas in location: %s; gasConsumed: %d, gasLimit: %d",
					ex.Descriptor, gasConsumed, gasLimit,
				)
				err = fmt.Errorf("out of gas: %s", log)
			default:
				panic(r)
			}
		}
	}()

	// After transaction execution, sync the gas consumed from GnoVM to SDK context
	newCtx, err = next(newCtx, tx, simulate)
	if err == nil {
		gasConsumed := gnoCtx.GasMeter().GasConsumed()
		ctx.GasMeter().ConsumeGas(storetypes.Gas(gasConsumed), "gnovm execution")
	}

	return newCtx, err
}
