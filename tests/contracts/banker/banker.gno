package banker

import (
	"chain"
	"chain/runtime"
	"chain/banker"
)

type activity struct {
	caller address
	sent   chain.Coins
}

var activities []*activity

// Deposit accepts coins sent with the transaction.
// The coins are automatically transferred to the realm's address.
func Deposit(_ realm) string {
	caller := runtime.PreviousRealm().Address()
	sent := banker.OriginSend()

	// Record the activity
	act := &activity{
		caller: caller,
		sent:   sent,
	}
	activities = append(activities, act)

	return "deposit recorded"
}

// GetBalance returns the balance of the realm itself.
func GetBalance(_ realm) chain.Coins {
	b := banker.NewBanker(banker.BankerTypeReadonly)
	realmAddr := runtime.CurrentRealm().Address()
	return b.GetCoins(realmAddr)
}

// GetUserBalance returns the balance of a specific user address.
func GetUserBalance(_ realm, addr address) chain.Coins {
	b := banker.NewBanker(banker.BankerTypeReadonly)
	return b.GetCoins(addr)
}

// SendCoins sends coins from the realm to a specified address.
func SendCoins(_ realm, to address, amount int64, denom string) string {
	b := banker.NewBanker(banker.BankerTypeOriginSend)
	realmAddr := runtime.CurrentRealm().Address()

	// Create the coins to send
	coins := chain.Coins{chain.Coin{Denom: denom, Amount: amount}}

	// Send the coins from realm to recipient
	b.SendCoins(realmAddr, to, coins)

	return "coins sent"
}

// Render displays the realm's current balance and recent activity.
func Render(_ string) string {
	b := banker.NewBanker(banker.BankerTypeReadonly)
	realmAddr := runtime.CurrentRealm().Address()
	balance := b.GetCoins(realmAddr)

	res := "# Banker Realm\n\n"
	res += "## Current Balance\n"
	res += balance.String() + "\n\n"

	res += "## Recent Activity\n"
	if len(activities) == 0 {
		res += "No activity yet\n"
	} else {
		// Show last 10 activities
		start := 0
		if len(activities) > 10 {
			start = len(activities) - 10
		}
		for i := start; i < len(activities); i++ {
			act := activities[i]
			res += "- " + act.caller.String() + " deposited " + act.sent.String() + "\n"
		}
	}

	return res
}
