package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewMsgRun creates a new MsgRun instance.
func NewMsgRun(caller string, send, maxDeposit sdk.Coins, pkg []byte) *MsgRun {
	return &MsgRun{
		Caller:     caller,
		Send:       send,
		MaxDeposit: maxDeposit,
		Pkg:        pkg,
	}
}
