package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RealmStorage returns the storage information for a given realm.
func (q queryServer) RealmStorage(ctx context.Context, req *types.QueryRealmStorageRequest) (*types.QueryRealmStorageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.PkgPath == "" {
		return nil, status.Error(codes.InvalidArgument, "package path cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	gnoCtx, err := q.k.BuildGnoContext(sdkCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize VM: %w", err)
	}

	result, err := q.k.VMKeeper.QueryStorage(gnoCtx, req.PkgPath)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to query realm storage")
	}

	return &types.QueryRealmStorageResponse{
		Result: result,
	}, nil
}
