package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
	gnoCtx, err := k.BuildGnoContext(sdkCtx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to initialize VM")
	}

	send := types.StdCoinsFromSDKCoins(msg.Send)
	maxDep := types.StdCoinsFromSDKCoins(msg.MaxDeposit)

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

	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case storetypes.ErrorOutOfGas:
				log := fmt.Sprintf(
					"out of gas from VM usage in location: %v; gasUsed: %d",
					rType.Descriptor, sdkCtx.GasMeter().GasConsumed())

				err = errorsmod.Wrap(sdkerrors.ErrOutOfGas, log)
			default:
				err = fmt.Errorf("panic while calling VM: %v (%v)", r, rType)
			}
		} else {
			// this commits the changes to the module store (that is only committed later)
			k.VMKeeper.CommitGnoTransactionStore(gnoCtx)
		}
	}()

	if err := k.VMKeeper.AddPackage(gnoCtx, vmMsg); err != nil {
		return nil, errorsmod.Wrap(err, "failed to add package")
	}

	return &types.MsgAddPackageResponse{}, nil
}
