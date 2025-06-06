package squads

import (
	"context"
	"squads/generated/squads_multisig_program"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type SQuard struct {
	programID solana.PublicKey
}

func NewSQuard() *SQuard {
	// devnet and mainnet-beta
	programID := solana.MustPublicKeyFromBase58("SQDS4ep65T869zMMBKyuUq6aD6EgTu8psMjkvj52pCf")
	return &SQuard{programID: programID}
}

func (s *SQuard) Account(ctx context.Context, client *rpc.Client, multisigPda solana.PublicKey) (*squads_multisig_program.Multisig, error) {
	out, err := client.GetAccountInfo(ctx, multisigPda)
	if err != nil {
		return nil, err
	}
	data := out.Value.Data.GetBinary()

	multisig := &squads_multisig_program.Multisig{}
	decoder := ag_binary.NewBorshDecoder(data)
	if err := multisig.UnmarshalWithDecoder(decoder); err != nil {
		return nil, err
	}
	return multisig, nil
}
