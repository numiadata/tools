package crypto

import (
	"reflect"

	amino "github.com/tendermint/go-amino"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	Ed25519PubKeyAminoName  = "tendermint/PubKeyEd25519"
	Ed25519PrivKeyAminoName = "tendermint/PrivKeyEd25519"

	SecpPrivKeyAminoName = "tendermint/PrivKeySecp256k1"
	SecpPubKeyAminoName  = "tendermint/PubKeySecp256k1"

	PubKeyMultisigThresholdAminoRoute = "tendermint/PubKeyMultisigThreshold"
)

var cdc = amino.NewCodec()

// nameTable is used to map public key concrete types back
// to their registered amino names. This should eventually be handled
// by amino. Example usage:
// nameTable[reflect.TypeOf(ed25519.PubKey{})] = ed25519.PubKeyAminoName
var nameTable = make(map[reflect.Type]string, 3)

func init() {
	// NOTE: It's important that there be no conflicts here,
	// as that would change the canonical representations,
	// and therefore change the address.
	// TODO: Remove above note when
	// https://github.com/tendermint/go-amino/issues/9
	// is resolved
	RegisterAmino(cdc)
}

// RegisterAmino registers all crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	// These are all written here instead of
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKey{},
		Ed25519PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PubKey{},
		SecpPubKeyAminoName, nil)
	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PrivKey{},
		Ed25519PrivKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKey{},
		SecpPrivKeyAminoName, nil)
}
