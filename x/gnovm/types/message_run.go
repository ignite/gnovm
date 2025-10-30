package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewMsgRun creates a new MsgRun instance.
func NewMsgRun(caller string, send sdk.Coins, maxDeposit sdk.Coin, pkg *Package) *MsgRun {
	return &MsgRun{
		Caller:     caller,
		Send:       send,
		MaxDeposit: maxDeposit,
		Pkg:        pkg,
	}
}

func (msg MsgRun) ValidateBasic() error {
	// TODO add other checks
	return msg.Pkg.Validate()
}
