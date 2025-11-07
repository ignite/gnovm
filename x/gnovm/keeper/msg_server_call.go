package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"

	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) Call(ctx context.Context, msg *types.MsgCall) (resp *types.MsgCallResponse, err error) {
	callerBytes, err := k.addressCodec.StringToBytes(msg.Caller)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to convert caller address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	gnoCtx, err := k.BuildGnoContext(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}

	vmMsg := vm.MsgCall{
		Caller:     types.ToCryptoAddress(callerBytes),
		Send:       types.StdCoinsFromSDKCoins(msg.Send),
		MaxDeposit: types.StdCoinsFromSDKCoins(msg.MaxDeposit),
		PkgPath:    msg.PkgPath,
		Func:       msg.Function,
		Args:       msg.Args,
	}

	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case storetypes.ErrorOutOfGas:
				log := fmt.Sprintf(
					"out of gas from VM usage in location: %v; gasUsed: %d",
					rType.Descriptor, sdkCtx.GasMeter().GasConsumed())

				err = errorsmod.Wrap(sdkerrors.ErrOutOfGas, log)
			default:
				err = fmt.Errorf("panic while calling VM: %v (%v)", r, rType)
			}
		} else {
			// this commits the changes to the module store (that is only committed later)
			k.VMKeeper.CommitGnoTransactionStore(gnoCtx)
		}
	}()

	result, err := k.VMKeeper.Call(
		gnoCtx,
		vmMsg,
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to call VM")
	}

	return &types.MsgCallResponse{
		Result: result,
	}, nil
}
