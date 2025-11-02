package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewMsgCall creates a new MsgCall instance.
func NewMsgCall(caller string, send, maxDeposit sdk.Coins, pkgPath string, function string, args []string) *MsgCall {
	return &MsgCall{
		Caller:     caller,
		Send:       send,
		MaxDeposit: maxDeposit,
		PkgPath:    pkgPath,
		Function:   function,
		Args:       args,
	}
}
