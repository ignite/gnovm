package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gnolang/gno/tm2/pkg/std"
)

// StdCoinsFromSDKCoins converts sdk.Coins to std.Coins
func StdCoinsFromSDKCoins(coins sdk.Coins) std.Coins {
	stdCoins := make(std.Coins, len(coins))
	for i, coin := range coins {
		stdCoins[i] = std.NewCoin(coin.Denom, coin.Amount.Int64())
	}
	return stdCoins
}

// SDKCoinsFromStdCoins converts std.Coins to sdk.Coins
func SDKCoinsFromStdCoins(amt std.Coins) sdk.Coins {
	coins := make(sdk.Coins, len(amt))
	for i, coin := range amt {
		coins[i] = sdk.NewInt64Coin(coin.Denom, coin.Amount)
	}
	return coins
}
