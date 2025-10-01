package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgRun(caller string, send sdk.Coins, pkg []byte) *MsgRun {
	return &MsgRun{
		Caller: caller,
		Send:   send,
		Pkg:    pkg,
	}
}
