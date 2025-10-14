package keeper

import (
	"context"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/gnovm/x/gnovm/types"
)

func (k msgServer) Run(ctx context.Context, msg *types.MsgRun) (*types.MsgRunResponse, error) {
	callerBytes, err := k.addressCodec.StringToBytes(msg.Caller)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to convert caller address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	gnoCtx, err := k.BuildGnoContextWithStore(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}
	defer k.VMKeeper.CommitGnoTransactionStore(gnoCtx)

	send := types.StdCoinsFromSDKCoins(msg.Send)
	maxDep := types.StdCoinsFromSDKCoins(sdk.NewCoins(msg.MaxDeposit))

	var mpkg std.MemPackage
	if err := json.Unmarshal(msg.Pkg, &mpkg); err != nil {
		return nil, errorsmod.Wrap(err, "invalid package")
	}
	if err := mpkg.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "invalid package")
	}

	if _, err := k.VMKeeper.Run(
		gnoCtx,
		vm.MsgRun{
			Caller:     types.ToCryptoAddress(callerBytes),
			Send:       send,
			MaxDeposit: maxDep,
			Package:    &mpkg,
		},
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to run VM")
	}

	return &types.MsgRunResponse{}, nil
}
