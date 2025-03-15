package keeper

import (
	"context"

	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/sdk/params"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/gnovm/x/gnovm/types"
)

// vmAuthKeeper is a wrapper of the Cosmos SDK auth keeper to the VM expected auth keeper.
type vmAuthKeeper struct {
	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper
}

// GetAccount implements vm.AccountKeeperI.
func (v vmAuthKeeper) GetAccount(ctx gnosdk.Context, addr crypto.Address) std.Account {
	account := v.authKeeper.GetAccount(context.TODO(), addr.Bytes())
	return types.StdAccountFromSDKAccount(account, v.bankKeeper)
}

// vmBankKeeper is a wrapper of the Cosmos SDK bank keeper to the VM expected bank keeper.
type vmBankKeeper struct {
	bankKeeper types.BankKeeper
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

var _ vm.ParamsKeeperI = (*vmKeeperParams)(nil)

type vmKeeperParams struct {
	k *Keeper
}

// GetAny implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetAny(ctx gnosdk.Context, key string) interface{} {
	panic("unimplemented")
}

// GetBool implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetBool(ctx gnosdk.Context, key string, ptr *bool) {
	panic("unimplemented")
}

// GetBytes implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetBytes(ctx gnosdk.Context, key string, ptr *[]byte) {
	panic("unimplemented")
}

// GetInt64 implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetInt64(ctx gnosdk.Context, key string, ptr *int64) {
	panic("unimplemented")
}

// GetRaw implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetRaw(ctx gnosdk.Context, key string) []byte {
	panic("unimplemented")
}

// GetString implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetString(ctx gnosdk.Context, key string, ptr *string) {
	panic("unimplemented")
}

// GetStrings implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetStrings(ctx gnosdk.Context, key string, ptr *[]string) {
	panic("unimplemented")
}

// GetStruct implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetStruct(ctx gnosdk.Context, key string, strctPtr interface{}) {
	panic("unimplemented")
}

// GetUint64 implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetUint64(ctx gnosdk.Context, key string, ptr *uint64) {
	panic("unimplemented")
}

// Has implements vm.ParamsKeeperI.
func (k *vmKeeperParams) Has(ctx gnosdk.Context, key string) bool {
	panic("unimplemented")
}

// SetAny implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetAny(ctx gnosdk.Context, key string, value interface{}) {
	panic("unimplemented")
}

// SetBool implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetBool(ctx gnosdk.Context, key string, value bool) {
	panic("unimplemented")
}

// SetBytes implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetBytes(ctx gnosdk.Context, key string, value []byte) {
	panic("unimplemented")
}

// SetInt64 implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetInt64(ctx gnosdk.Context, key string, value int64) {
	panic("unimplemented")
}

// SetRaw implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetRaw(ctx gnosdk.Context, key string, value []byte) {
	panic("unimplemented")
}

// SetString implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetString(ctx gnosdk.Context, key string, value string) {
	panic("unimplemented")
}

// SetStrings implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetStrings(ctx gnosdk.Context, key string, value []string) {
	panic("unimplemented")
}

// SetStruct implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetStruct(ctx gnosdk.Context, key string, strct interface{}) {
	panic("unimplemented")
}

// SetUint64 implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetUint64(ctx gnosdk.Context, key string, value uint64) {
	panic("unimplemented")
}

// IsRegistered implements vm.ParamsKeeperI.
func (k *vmKeeperParams) IsRegistered(moduleName string) bool {
	panic("unimplemented")
}

// GetRegisteredKeeper implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetRegisteredKeeper(moduleName string) params.ParamfulKeeper {
	panic("unimplemented")
}
