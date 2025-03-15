package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) AddPackage(ctx context.Context, msg *types.MsgAddPackage) (*types.MsgAddPackageResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgAddPackageResponse{}, nil
}
