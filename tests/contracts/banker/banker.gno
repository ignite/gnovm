package banker

import (
	"std"
)

type activity struct {
	caller std.Address
	sent   std.Coins
	time   int64
}

var activities []*activity

// Deposit accepts coins sent with the transaction.
// The coins are automatically transferred to the realm's address.
func Deposit(_ realm) string {
	caller := std.PrevRealm().Addr()
	sent := std.GetOrigSend()

	// Record the activity
	act := &activity{
		caller: caller,
		sent:   sent,
		time:   std.GetTimestamp(),
	}
	activities = append(activities, act)

	return "deposit recorded"
}

// GetBalance returns the balance of the realm itself.
func GetBalance(_ realm) std.Coins {
	banker := std.GetBanker(std.BankerTypeReadonly)
	realmAddr := std.CurrentRealm().Addr()
	return banker.GetCoins(realmAddr)
}

// GetUserBalance returns the balance of a specific user address.
func GetUserBalance(_ realm, addr std.Address) std.Coins {
	banker := std.GetBanker(std.BankerTypeReadonly)
	return banker.GetCoins(addr)
}

// SendCoins sends coins from the realm to a specified address.
func SendCoins(_ realm, to std.Address, amount int64, denom string) string {
	banker := std.GetBanker(std.BankerTypeOrigSend)
	realmAddr := std.CurrentRealm().Addr()

	// Create the coins to send
	coins := std.Coins{std.Coin{Denom: denom, Amount: amount}}

	// Send the coins from realm to recipient
	banker.SendCoins(realmAddr, to, coins)

	return "coins sent"
}

// Render displays the realm's current balance and recent activity.
func Render(_ string) string {
	banker := std.GetBanker(std.BankerTypeReadonly)
	realmAddr := std.CurrentRealm().Addr()
	balance := banker.GetCoins(realmAddr)

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
