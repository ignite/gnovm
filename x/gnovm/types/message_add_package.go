package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgAddPackage(creator string, send, maxDeposit sdk.Coins, pkg []byte) *MsgAddPackage {
	return &MsgAddPackage{
		Creator:    creator,
		Send:       send,
		MaxDeposit: maxDeposit,
		Package:    pkg,
	}
}
