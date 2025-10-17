package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/types"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) Render(ctx context.Context, req *types.QueryRenderRequest) (*httpbody.HttpBody, error) {
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

	// Construct render expression
	var expr string
	if len(req.Args) > 0 {
		expr = fmt.Sprintf(`Render(%q)`, strings.Join(req.GetArgs(), `","`))
	} else {
		expr = `Render("")`
	}

	result, err := q.k.VMKeeper.QueryEval(gnoCtx, req.PkgPath, expr)
	if err != nil {
		return nil, fmt.Errorf("failed to render contract: %w", err)
	}

	return &httpbody.HttpBody{
		ContentType: "text/html",
		Data:        []byte(result),
	}, nil
}
