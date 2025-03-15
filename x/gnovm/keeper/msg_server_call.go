package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) Call(ctx context.Context, msg *types.MsgCall) (*types.MsgCallResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Caller); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgCallResponse{}, nil
}
