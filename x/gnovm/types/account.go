package types

import (
	"context"

	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/std"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StfAccountFromSDKAccount convers an sdk.AccountI as an std.Account
// As std.Account have as well balances, it takes a bank keeper as well.
func StdAccountFromSDKAccount(ctx context.Context, acc sdk.AccountI, bankKeeper BankKeeper) std.Account {
	return &accountWrapper{
		ctx:        ctx,
		acc:        acc,
		bankKeeper: bankKeeper,
	}
}

var _ std.Account = (*accountWrapper)(nil)

type accountWrapper struct {
	ctx        context.Context
	acc        sdk.AccountI
	bankKeeper BankKeeper
}

// GetAccountNumber implements std.Account.
func (a *accountWrapper) GetAccountNumber() uint64 {
	return a.acc.GetAccountNumber()
}

// GetAddress implements std.Account.
func (a *accountWrapper) GetAddress() crypto.Address {
	return crypto.AddressFromBytes(a.acc.GetAddress())
}

// GetCoins implements std.Account.
func (a *accountWrapper) GetCoins() std.Coins {
	coins := a.bankKeeper.GetAllBalances(a.ctx, a.acc.GetAddress())
	return StdCoinsFromSDKCoins(coins)
}

// GetPubKey implements std.Account.
func (a *accountWrapper) GetPubKey() crypto.PubKey {
	return PubKeyFromSDKPubKey(a.acc.GetPubKey())
}

// GetSequence implements std.Account.
func (a *accountWrapper) GetSequence() uint64 {
	return a.acc.GetSequence()
}

// SetAccountNumber implements std.Account.
func (a *accountWrapper) SetAccountNumber(accountNumber uint64) error {
	return a.acc.SetAccountNumber(accountNumber)
}

// SetAddress implements std.Account.
func (a *accountWrapper) SetAddress(addr crypto.Address) error {
	return a.acc.SetAddress(addr.Bytes())
}

// SetSequence implements std.Account.
func (a *accountWrapper) SetSequence(seq uint64) error {
	return a.acc.SetSequence(seq)
}

// SetCoins implements std.Account.
// It is a no-op as the account handling is done by the x/bank module.
func (a *accountWrapper) SetCoins(coins std.Coins) error {
	return nil
}

// SetPubKey implements std.Account.
// It is a no-op as the account handling is done by the x/auth module.
func (a *accountWrapper) SetPubKey(pub crypto.PubKey) error {
	return nil
}

// String implements std.Account.
func (a *accountWrapper) String() string {
	return a.acc.String()
}
