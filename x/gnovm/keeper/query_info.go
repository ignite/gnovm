package keeper

import (
	"context"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Info queries the package internal info by its package path.
func (q queryServer) Info(ctx context.Context, req *types.QueryInfoRequest) (*types.QueryInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.PkgPath == "" {
		return nil, status.Error(codes.InvalidArgument, "package path cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	gnoCtx, err := q.k.BuildGnoContext(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}

	result, err := q.k.VMKeeper.QueryDoc(gnoCtx, req.PkgPath)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to query package info")
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to marshal package info")
	}

	return &types.QueryInfoResponse{Result: string(jsonBytes)}, nil
}
