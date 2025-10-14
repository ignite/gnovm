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

func (k msgServer) AddPackage(ctx context.Context, msg *types.MsgAddPackage) (*types.MsgAddPackageResponse, error) {
	creatorBytes, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := (&k).initializeVMKeeper(sdkCtx); err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}

	gnoCtx, err := k.BuildGnoContextWithStore(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}
	defer k.VMKeeper.CommitGnoTransactionStore(gnoCtx)

	send := types.StdCoinsFromSDKCoins(msg.Deposit)
	maxDep := types.StdCoinsFromSDKCoins(sdk.NewCoins(msg.MaxDeposit))

	var mpkg std.MemPackage
	if err := json.Unmarshal(msg.Package, &mpkg); err != nil {
		return nil, errorsmod.Wrap(err, "invalid package")
	}
	if err := mpkg.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "invalid package")
	}

	vmMsg := vm.MsgAddPackage{
		Creator:    types.ToCryptoAddress(creatorBytes),
		Package:    &mpkg,
		Send:       send,
		MaxDeposit: maxDep,
	}
	if err := k.VMKeeper.AddPackage(gnoCtx, vmMsg); err != nil {
		return nil, errorsmod.Wrap(err, "failed to add package")
	}

	return &types.MsgAddPackageResponse{}, nil
}
