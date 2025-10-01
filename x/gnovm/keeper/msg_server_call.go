package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) Call(ctx context.Context, msg *types.MsgCall) (*types.MsgCallResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Caller); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	res, err := k.VMKeeper.Call(
		types.GnoContextFromSDKContext(sdkCtx),
		vm.MsgCall{}, // TODO: Handle the message
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to call VM")
	}

	return &types.MsgCallResponse{
		Result: res,
	}, nil
}
