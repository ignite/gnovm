package types

import "github.com/gnolang/gno/tm2/pkg/crypto"

// ToCryptoAddress converts a byte slice to crypto.Address safely.
// If the input length is not 20 bytes, it returns the zero address.
func ToCryptoAddress(b []byte) crypto.Address {
	var addr crypto.Address
	if len(b) == len(addr) {
		copy(addr[:], b)
	}
	return addr
}
