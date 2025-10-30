package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgAddPackage(creator string, deposit sdk.Coins, maxDeposit sdk.Coin, pkg *Package) *MsgAddPackage {
	return &MsgAddPackage{
		Creator:    creator,
		Deposit:    deposit,
		MaxDeposit: maxDeposit,
		Package:    pkg,
	}
}

func (msg MsgAddPackage) ValidateBasic() error {
	// TODO add other checks
	return msg.Package.Validate()
}
