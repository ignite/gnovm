package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Eval(ctx context.Context, req *types.QueryEvalRequest) (*types.QueryEvalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.PkgPath == "" {
		return nil, status.Error(codes.InvalidArgument, "package path cannot be empty")
	}

	if req.Expr == "" {
		return nil, status.Error(codes.InvalidArgument, "expression cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	gnoCtx, err := q.k.BuildGnoContext(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}

	result, err := q.k.VMKeeper.QueryEval(gnoCtx, req.PkgPath, req.Expr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to query package eval")
	}

	return &types.QueryEvalResponse{Result: result}, nil
}
