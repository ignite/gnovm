package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) Call(ctx context.Context, msg *types.MsgCall) (*types.MsgCallResponse, error) {
	callerBytes, err := k.addressCodec.StringToBytes(msg.Caller)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to convert caller address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	gnoCtx, err := k.BuildGnoContextWithStore(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}
	defer k.VMKeeper.CommitGnoTransactionStore(gnoCtx)

	vmMsg := vm.MsgCall{
		Caller:     types.ToCryptoAddress(callerBytes),
		Send:       types.StdCoinsFromSDKCoins(msg.Send),
		MaxDeposit: types.StdCoinsFromSDKCoins(sdk.NewCoins(msg.MaxDeposit)),
		PkgPath:    msg.PkgPath,
		Func:       msg.Function,
		Args:       msg.Args,
	}

	var result string
	var callErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				callErr = fmt.Errorf("panic while calling VM: %v", r)
			}
		}()
		result, callErr = k.VMKeeper.Call(
			gnoCtx,
			vmMsg,
		)
	}()
	if callErr != nil {
		return nil, errorsmod.Wrap(callErr, "failed to call VM")
	}

	return &types.MsgCallResponse{
		Result: result,
	}, nil
}
