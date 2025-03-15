package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgCall(caller string, send sdk.Coins, pkgPath string, function string, args []string) *MsgCall {
	return &MsgCall{
		Caller:   caller,
		Send:     send,
		PkgPath:  pkgPath,
		Function: function,
		Args:     args,
	}
}
