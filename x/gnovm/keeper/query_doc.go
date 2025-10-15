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

// Doc queries the documentation by its package path.
func (q queryServer) Doc(ctx context.Context, req *types.QueryDocRequest) (*types.QueryDocResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.PkgPath == "" {
		return nil, status.Error(codes.InvalidArgument, "package path cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	gnoCtx, err := q.k.BuildGnoContextWithStore(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}

	result, err := q.k.VMKeeper.QueryDoc(gnoCtx, req.PkgPath)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to query package documentation")
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to marshal package documentation")
	}

	return &types.QueryDocResponse{Content: string(jsonBytes)}, nil
}
