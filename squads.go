package squads

import (
	"context"

	"github.com/Lee0x273/go-squads/generated/squads_multisig_program"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// CreateMultisigIx return a "create multisig wallet" instruction
// createKey:One-time, used to generate multisigPda
// creator: signer and feePayer
func CreateMultisigIx(ctx context.Context, client *rpc.Client, createKey, creator solana.PublicKey, members []squads_multisig_program.Member, threshold uint16, timelock uint32) (solana.Instruction, solana.PublicKey, error) {
	args := squads_multisig_program.MultisigCreateArgsV2{
		ConfigAuthority: nil,
		Threshold:       threshold,
		Members:         members,
		TimeLock:        timelock,
		RentCollector:   nil,
	}
	programConfigPda, _ := GetProgramConfigPda()
	var programConfig squads_multisig_program.ProgramConfig
	if err := client.GetAccountDataInto(ctx, programConfigPda, &programConfig); err != nil {
		return nil, solana.PublicKey{}, err
	}
	var configTreasury = programConfig.Treasury
	multisigPda, _ := GetMultisigPda(createKey)

	ix := squads_multisig_program.NewMultisigCreateV2Instruction(
		args,
		programConfigPda,
		configTreasury,
		multisigPda, // creator
		createKey,
		creator,
		solana.SystemProgramID,
	).Build()
	return ix, multisigPda, nil
}
