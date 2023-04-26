package state

import (
	amino "github.com/tendermint/go-amino"

	cryptoamino "github.com/numiadata/tools/cosi/utils/amino/crypto"
)

var cdc = amino.NewCodec()

func init() {
	cryptoamino.RegisterAmino(cdc)
}
