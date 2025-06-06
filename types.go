package squads

import (
	"github.com/gagliardetto/solana-go"
)

var SYSTEM_PROGRAM_ID = solana.MustPublicKeyFromBase58("11111111111111111111111111111111")

const PROGRAM_ADDRESS = "SQDS4ep65T869zMMBKyuUq6aD6EgTu8psMjkvj52pCf"

var PROGRAM_ID = solana.MustPublicKeyFromBase58(PROGRAM_ADDRESS)

type AddressLookupTableState struct {
	DeactivationSlot           uint64
	LastExtendedSlot           uint64
	LastExtendedSlotStartIndex uint8
	Addresses                  []solana.PublicKey
}

type AddressLookupTableAccount struct {
	Key   solana.PublicKey
	State AddressLookupTableState
}

type Permission uint8

const (
	Initiate Permission = 1 << 0
	Vote     Permission = 1 << 1
	Execute  Permission = 1 << 2
)

func (p Permission) Has(permission Permission) bool {
	return p&permission != 0
}
