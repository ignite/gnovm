package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"

	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) Run(ctx context.Context, msg *types.MsgRun) (*types.MsgRunResponse, error) {
	callerBytes, err := k.addressCodec.StringToBytes(msg.Caller)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to convert caller address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	gnoCtx, err := k.BuildGnoContext(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}

	send := types.StdCoinsFromSDKCoins(msg.Send)
	maxDep := types.StdCoinsFromSDKCoins(sdk.NewCoins(msg.MaxDeposit))

	if _, err := k.VMKeeper.Run(
		gnoCtx,
		vm.MsgRun{
			Caller:     types.ToCryptoAddress(callerBytes),
			Send:       send,
			MaxDeposit: maxDep,
			Package:    msg.Pkg.ToMemPackage(),
		},
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to run VM")
	}

	// this commits the changes to the module store (that is only committed later)
	k.VMKeeper.CommitGnoTransactionStore(gnoCtx)

	return &types.MsgRunResponse{}, nil
}
