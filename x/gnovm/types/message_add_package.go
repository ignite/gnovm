package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgAddPackage(creator string, deposit sdk.Coins, maxDeposit sdk.Coin, pkg []byte) *MsgAddPackage {
	return &MsgAddPackage{
		Creator:    creator,
		Deposit:    deposit,
		MaxDeposit: maxDeposit,
		Package:    pkg,
	}
}
