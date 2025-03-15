package keeper

import (
	"context"

	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/gnovm/x/gnovm/types"
)

type vmAuthKeeper struct {
	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper
}

// NewVMAuthKeeper is a wrapper of the Cosmos SDK auth keeper to the VM expected auth keeper.
func NewVMAuthKeeper(authKeeper types.AuthKeeper, bankKeeper types.BankKeeper) vm.AccountKeeperI {
	return vmAuthKeeper{authKeeper: authKeeper, bankKeeper: bankKeeper}
}

// GetAccount implements vm.AccountKeeperI.
func (v vmAuthKeeper) GetAccount(ctx gnosdk.Context, addr crypto.Address) std.Account {
	account := v.authKeeper.GetAccount(context.TODO(), addr.Bytes())
	return types.StdAccountFromSDKAccount(account, v.bankKeeper)
}

type vmBankKeeper struct {
	bankKeeper types.BankKeeper
}

// NewVMBankKeeper is a wrapper of the Cosmos SDK bank keeper to the VM expected bank keeper.
func NewVMBankKeeper(k types.BankKeeper) vm.BankKeeperI {
	return vmBankKeeper{bankKeeper: k}
}

// AddCoins implements vm.BankKeeperI.
func (v vmBankKeeper) AddCoins(ctx gnosdk.Context, addr crypto.Address, amt std.Coins) (std.Coins, error) {
	panic("unimplemented")
}

// GetCoins implements vm.BankKeeperI.
func (v vmBankKeeper) GetCoins(ctx gnosdk.Context, addr crypto.Address) std.Coins {
	coins := v.bankKeeper.GetAllBalances(context.TODO(), addr.Bytes())
	return types.StdCoinsFromSDKCoins(coins)
}

// SendCoins implements vm.BankKeeperI.
func (v vmBankKeeper) SendCoins(ctx gnosdk.Context, fromAddr crypto.Address, toAddr crypto.Address, amt std.Coins) error {
	return v.bankKeeper.SendCoins(
		context.TODO(),
		fromAddr.Bytes(),
		toAddr.Bytes(),
		types.SDKCoinsFromStdCoins(amt),
	)
}

// SubtractCoins implements vm.BankKeeperI.
func (v vmBankKeeper) SubtractCoins(ctx gnosdk.Context, addr crypto.Address, amt std.Coins) (std.Coins, error) {
	panic("unimplemented")
}
