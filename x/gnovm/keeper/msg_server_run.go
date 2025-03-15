package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) Run(ctx context.Context, msg *types.MsgRun) (*types.MsgRunResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Caller); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgRunResponse{}, nil
}
