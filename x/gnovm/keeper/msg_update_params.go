package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	authority, err := k.addressCodec.StringToBytes(req.Authority)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", expectedAuthorityStr, req.Authority)
	}

	// Treat zero-valued params as no-op (skip validation and update)
	if req.Params.SysnamesPkgpath == "" &&
		req.Params.ChainDomain == "" &&
		req.Params.DefaultDeposit == "" &&
		req.Params.StoragePrice == "" &&
		len(req.Params.StorageFeeCollector) == 0 {
		return &types.MsgUpdateParamsResponse{}, nil
	}

	if err := req.Params.Validate(); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	gnoCtx, err := k.BuildGnoContext(sdkCtx)
	if err != nil {
		return nil, err
	}

	if err := k.VMKeeper.SetParams(gnoCtx, req.Params.ToVmParams()); err != nil {
		return nil, err
	}

	// duplicate to SDK module state
	if err := k.Params.Set(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
