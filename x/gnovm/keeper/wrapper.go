package keeper

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/sdk/params"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/gnovm/x/gnovm/types"
)

var _ vm.AccountKeeperI = (*vmAuthKeeper)(nil)

// vmAuthKeeper is a wrapper of the Cosmos SDK auth keeper to the VM expected auth keeper.
type vmAuthKeeper struct {
	logger     log.Logger
	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper
}

// GetAccount implements vm.AccountKeeperI.
func (v vmAuthKeeper) GetAccount(ctx gnosdk.Context, addr crypto.Address) std.Account {
	account := v.authKeeper.GetAccount(types.SDKContextFromGnoContext(ctx, v.logger), addr.Bytes())
	return types.StdAccountFromSDKAccount(account, v.bankKeeper)
}

var _ vm.BankKeeperI = (*vmBankKeeper)(nil)

// vmBankKeeper is a wrapper of the Cosmos SDK bank keeper to the VM expected bank keeper.
type vmBankKeeper struct {
	logger     log.Logger
	bankKeeper types.BankKeeper
}

// RestrictedDenoms implements vm.BankKeeperI.
func (v vmBankKeeper) RestrictedDenoms(ctx gnosdk.Context) []string {
	return []string{}
}

// AddCoins implements vm.BankKeeperI.
func (v vmBankKeeper) AddCoins(ctx gnosdk.Context, addr crypto.Address, amt std.Coins) (std.Coins, error) {
	panic("unimplemented")
}

// GetCoins implements vm.BankKeeperI.
func (v vmBankKeeper) GetCoins(ctx gnosdk.Context, addr crypto.Address) std.Coins {
	coins := v.bankKeeper.GetAllBalances(types.SDKContextFromGnoContext(ctx, v.logger), addr.Bytes())
	return types.StdCoinsFromSDKCoins(coins)
}

// SendCoins implements vm.BankKeeperI.
func (v vmBankKeeper) SendCoins(ctx gnosdk.Context, fromAddr crypto.Address, toAddr crypto.Address, amt std.Coins) error {
	return v.bankKeeper.SendCoins(
		types.SDKContextFromGnoContext(ctx, v.logger),
		fromAddr.Bytes(),
		toAddr.Bytes(),
		types.SDKCoinsFromStdCoins(amt),
	)
}

// SendCoinsUnrestricted implements vm.BankKeeperI.
func (v vmBankKeeper) SendCoinsUnrestricted(ctx gnosdk.Context, fromAddr crypto.Address, toAddr crypto.Address, amt std.Coins) error {
	return v.bankKeeper.SendCoins(
		types.SDKContextFromGnoContext(ctx, v.logger),
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
	k      *Keeper
	sdkCtx sdk.Context
}

// SetSDKContext sets the SDK context for store operations
func (k *vmKeeperParams) SetSDKContext(ctx sdk.Context) {
	k.sdkCtx = ctx
}

// paramStoreKey generates the store key for a given parameter
func (k *vmKeeperParams) paramStoreKey(key string) string {
	return fmt.Sprintf("params/%s", key)
}

// GetAny implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetAny(ctx gnosdk.Context, key string) interface{} {
	// get raw value from the store
	data := k.GetRaw(ctx, key)
	if len(data) == 0 {
		return nil
	}

	// unmarshal the value
	var value interface{}
	err := json.Unmarshal(data, &value)
	if err != nil {
		k.k.logger.Error("failed to unmarshal param value", "key", key, "error", err)
		return nil
	}

	return value
}

// GetBool implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetBool(ctx gnosdk.Context, key string, ptr *bool) {
	// get raw value from the store
	data := k.GetRaw(ctx, key)
	if len(data) == 0 {
		*ptr = false
		return
	}

	// unmarshal the value
	var value bool
	err := json.Unmarshal(data, &value)
	if err != nil {
		k.k.logger.Error("failed to unmarshal bool param value", "key", key, "error", err)
		*ptr = false
		return
	}

	*ptr = value
}

// GetBytes implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetBytes(ctx gnosdk.Context, key string, ptr *[]byte) {
	// get raw value from the store
	data := k.GetRaw(ctx, key)
	if len(data) == 0 {
		*ptr = []byte{}
		return
	}

	*ptr = data
}

// GetInt64 implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetInt64(ctx gnosdk.Context, key string, ptr *int64) {
	// get raw value from the store
	data := k.GetRaw(ctx, key)
	if len(data) == 0 {
		*ptr = 0
		return
	}

	// unmarshal the value
	var value int64
	err := json.Unmarshal(data, &value)
	if err != nil {
		k.k.logger.Error("failed to unmarshal int64 param value", "key", key, "error", err)
		*ptr = 0
		return
	}

	*ptr = value
}

// GetRaw implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetRaw(ctx gnosdk.Context, key string) []byte {
	// use store service to get the value
	store := k.k.storeService.OpenKVStore(k.sdkCtx)
	storeKey := []byte(k.paramStoreKey(key))

	data, err := store.Get(storeKey)
	if err != nil {
		k.k.logger.Error("failed to get param value from store", "key", key, "error", err)
		return []byte{}
	}

	return data
}

// GetString implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetString(ctx gnosdk.Context, key string, ptr *string) {
	// get raw value from the store
	data := k.GetRaw(ctx, key)
	if len(data) == 0 {
		*ptr = ""
		return
	}

	// unmarshal the value
	var value string
	err := json.Unmarshal(data, &value)
	if err != nil {
		k.k.logger.Error("failed to unmarshal string param value", "key", key, "error", err)
		*ptr = ""
		return
	}

	*ptr = value
}

// GetStrings implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetStrings(ctx gnosdk.Context, key string, ptr *[]string) {
	// get raw value from the store
	data := k.GetRaw(ctx, key)
	if len(data) == 0 {
		*ptr = []string{}
		return
	}

	// unmarshal the value
	var value []string
	err := json.Unmarshal(data, &value)
	if err != nil {
		k.k.logger.Error("failed to unmarshal strings param value", "key", key, "error", err)
		*ptr = []string{}
		return
	}

	*ptr = value
}

// GetStruct implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetStruct(ctx gnosdk.Context, key string, strctPtr interface{}) {
	// get raw value from the store
	data := k.GetRaw(ctx, key)
	if len(data) == 0 {
		return
	}

	// unmarshal the value
	err := json.Unmarshal(data, strctPtr)
	if err != nil {
		k.k.logger.Error("failed to unmarshal struct param value", "key", key, "error", err)
		return
	}
}

// GetUint64 implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetUint64(ctx gnosdk.Context, key string, ptr *uint64) {
	// get raw value from the store
	data := k.GetRaw(ctx, key)
	if len(data) == 0 {
		*ptr = 0
		return
	}

	// unmarshal the value
	var value uint64
	err := json.Unmarshal(data, &value)
	if err != nil {
		k.k.logger.Error("failed to unmarshal uint64 param value", "key", key, "error", err)
		*ptr = 0
		return
	}

	*ptr = value
}

// Has implements vm.ParamsKeeperI.
func (k *vmKeeperParams) Has(ctx gnosdk.Context, key string) bool {
	// use store service to check if the key exists
	store := k.k.storeService.OpenKVStore(k.sdkCtx)
	storeKey := []byte(k.paramStoreKey(key))

	has, err := store.Has(storeKey)
	if err != nil {
		k.k.logger.Error("failed to check if param exists", "key", key, "error", err)
		return false
	}

	return has
}

// SetAny implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetAny(ctx gnosdk.Context, key string, value interface{}) {
	// marshal the value
	data, err := json.Marshal(value)
	if err != nil {
		k.k.logger.Error("failed to marshal param value", "key", key, "error", err)
		return
	}

	// set the raw value
	k.SetRaw(ctx, key, data)
}

// SetBool implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetBool(ctx gnosdk.Context, key string, value bool) {
	// marshal the value
	data, err := json.Marshal(value)
	if err != nil {
		k.k.logger.Error("failed to marshal bool param value", "key", key, "error", err)
		return
	}

	// set the raw value
	k.SetRaw(ctx, key, data)
}

// SetBytes implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetBytes(ctx gnosdk.Context, key string, value []byte) {
	// directly use the raw bytes
	k.SetRaw(ctx, key, value)
}

// SetInt64 implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetInt64(ctx gnosdk.Context, key string, value int64) {
	// marshal the value
	data, err := json.Marshal(value)
	if err != nil {
		k.k.logger.Error("failed to marshal int64 param value", "key", key, "error", err)
		return
	}

	// set the raw value
	k.SetRaw(ctx, key, data)
}

// SetRaw implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetRaw(ctx gnosdk.Context, key string, value []byte) {
	// use store service to set the value
	store := k.k.storeService.OpenKVStore(k.sdkCtx)
	storeKey := []byte(k.paramStoreKey(key))

	if err := store.Set(storeKey, value); err != nil {
		k.k.logger.Error("failed to set param value in store", "key", key, "error", err)
	}
}

// SetString implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetString(ctx gnosdk.Context, key string, value string) {
	// marshal the value
	data, err := json.Marshal(value)
	if err != nil {
		k.k.logger.Error("failed to marshal string param value", "key", key, "error", err)
		return
	}

	// set the raw value
	k.SetRaw(ctx, key, data)
}

// SetStrings implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetStrings(ctx gnosdk.Context, key string, value []string) {
	// marshal the value
	data, err := json.Marshal(value)
	if err != nil {
		k.k.logger.Error("failed to marshal strings param value", "key", key, "error", err)
		return
	}

	// set the raw value
	k.SetRaw(ctx, key, data)
}

// SetStruct implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetStruct(ctx gnosdk.Context, key string, strct interface{}) {
	// marshal the value
	data, err := json.Marshal(strct)
	if err != nil {
		k.k.logger.Error("failed to marshal struct param value", "key", key, "error", err)
		return
	}

	// set the raw value
	k.SetRaw(ctx, key, data)
}

// SetUint64 implements vm.ParamsKeeperI.
func (k *vmKeeperParams) SetUint64(ctx gnosdk.Context, key string, value uint64) {
	// marshal the value
	data, err := json.Marshal(value)
	if err != nil {
		k.k.logger.Error("failed to marshal uint64 param value", "key", key, "error", err)
		return
	}

	// set the raw value
	k.SetRaw(ctx, key, data)
}

// IsRegistered implements vm.ParamsKeeperI.
func (k *vmKeeperParams) IsRegistered(moduleName string) bool {
	// all modules are considered registered for compatibility
	return true
}

// GetRegisteredKeeper implements vm.ParamsKeeperI.
func (k *vmKeeperParams) GetRegisteredKeeper(moduleName string) params.ParamfulKeeper {
	// return nil as we don't implement the full ParamfulKeeper interface
	// gnovm params are handled directly by our keeper
	return nil
}
