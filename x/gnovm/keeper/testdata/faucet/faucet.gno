package faucet

import (
	"chain"
	"chain/banker"
	"chain/runtime"
)

func Transfer(realm) {
	var (
		from = runtime.CurrentRealm().Address()
		to   = runtime.OriginCaller()
		amt  = chain.NewCoins(chain.NewCoin("ugnot", 1))
	)
	banker.NewBanker(banker.BankerTypeRealmSend).SendCoins(from, to, amt)
}
