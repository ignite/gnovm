package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	bft "github.com/gnolang/gno/tm2/pkg/bft/types"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) Call(ctx context.Context, msg *types.MsgCall) (*types.MsgCallResponse, error) {
	callerBytes, err := k.addressCodec.StringToBytes(msg.Caller)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to convert caller address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := (&k).initializeVMKeeper(sdkCtx); err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}

	gnoCtx := gnosdk.NewContext(
		gnosdk.RunTxModeDeliver,
		nil, // MultiStore provided by our wrapper
		&bft.Header{ChainID: sdkCtx.ChainID()},
		types.NewSlogFromCosmosLogger(k.logger),
	)

	vmMsg := vm.MsgCall{
		Caller:     types.ToCryptoAddress(callerBytes),
		Send:       types.StdCoinsFromSDKCoins(msg.Send),
		MaxDeposit: types.StdCoinsFromSDKCoins(sdk.NewCoins(msg.MaxDeposit)),
		PkgPath:    msg.PkgPath,
		Func:       msg.Function,
		Args:       msg.Args,
	}

	res, err := k.VMKeeper.Call(
		gnoCtx,
		vmMsg,
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to call VM")
	}

	return &types.MsgCallResponse{
		Result: res,
	}, nil
}
