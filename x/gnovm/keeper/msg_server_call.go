package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	var committed bool
	defer func() {
		if !committed && err == nil {
			// Commit the transaction store on successful execution
			k.VMKeeper.CommitGnoTransactionStore(gnoCtx)
			committed = true
		}
	}()

	vmMsg := vm.MsgCall{
		Caller:     types.ToCryptoAddress(callerBytes),
		Send:       types.StdCoinsFromSDKCoins(msg.Send),
		MaxDeposit: types.StdCoinsFromSDKCoins(sdk.NewCoins(msg.MaxDeposit)),
		PkgPath:    msg.PkgPath,
		Func:       msg.Function,
		Args:       msg.Args,
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic while calling VM: %v", r)
		}
	}()
	result, err := k.VMKeeper.Call(
		gnoCtx,
		vmMsg,
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to call VM")
	}

	// Commit the transaction store before returning success
	k.VMKeeper.CommitGnoTransactionStore(gnoCtx)
	committed = true

	return &types.MsgCallResponse{
		Result: result,
	}, nil
}
