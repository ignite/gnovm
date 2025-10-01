package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgAddPackage(creator string, deposit sdk.Coins) *MsgAddPackage {
	return &MsgAddPackage{
		Creator: creator,
		Deposit: deposit,
	}
}
