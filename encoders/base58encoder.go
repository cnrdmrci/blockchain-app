package encoders

import (
	"blockchain-app/handlers"
	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) string {
	return base58.Encode(input)
}

func Base58Decode(input string) []byte {
	decode, err := base58.Decode(input[:])
	handlers.HandleErrors(err)

	return decode
}
