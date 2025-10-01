package types

import (
	"bytes"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/gnolang/gno/tm2/pkg/crypto"
)

// PubKeyFromSDKPubKey converts a sdk crypto.PubKey to a gno crypto.PubKey
func PubKeyFromSDKPubKey(pub cryptotypes.PubKey) crypto.PubKey {
	return &gnoPubkeyWrapper{pub: pub}
}

var _ crypto.PubKey = (*gnoPubkeyWrapper)(nil)

type gnoPubkeyWrapper struct {
	pub cryptotypes.PubKey
}

// Address implements crypto.PubKey.
func (c *gnoPubkeyWrapper) Address() crypto.Address {
	return crypto.AddressFromBytes(c.pub.Address().Bytes())
}

// Bytes implements crypto.PubKey.
func (c *gnoPubkeyWrapper) Bytes() []byte {
	return c.pub.Bytes()
}

// Equals implements crypto.PubKey.
func (c *gnoPubkeyWrapper) Equals(pub crypto.PubKey) bool {
	return bytes.Equal(c.pub.Bytes(), pub.Bytes())
}

// String implements crypto.PubKey.
func (c *gnoPubkeyWrapper) String() string {
	return c.pub.String()
}

// VerifyBytes implements crypto.PubKey.
func (c *gnoPubkeyWrapper) VerifyBytes(msg []byte, sig []byte) bool {
	return c.pub.VerifySignature(msg, sig)
}
