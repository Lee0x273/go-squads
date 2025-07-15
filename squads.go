package squads

import (
	"context"

	"github.com/Lee0x273/go-squads/generated/squads_multisig_program"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// CreateMultisigIx returns a "create multisig wallet" instruction
// Parameters:
// - ctx: Context for the operation
// - client: RPC client for Solana
// - createKey: One-time key used to generate multisigPda
// - creator: Signer and fee payer
// - configAuthority: Optional configuration authority
// - members: List of members for the multisig
// - threshold: Number of signatures required for approval
// - timelock: Time lock period
// Returns:
// - Instruction for creating a multisig wallet
// - Public key of the created multisig
// - Error, if any
func CreateMultisigIx(ctx context.Context, client *rpc.Client, createKey, creator solana.PublicKey, configAuthority *solana.PublicKey, members []squads_multisig_program.Member, threshold uint16, timelock uint32, rentCollector *solana.PublicKey) (solana.Instruction, solana.PublicKey, error) {
	args := squads_multisig_program.MultisigCreateArgsV2{
		ConfigAuthority: configAuthority,
		Threshold:       threshold,
		Members:         members,
		TimeLock:        timelock,
		RentCollector:   rentCollector,
	}
	programConfigPda, err := GetProgramConfigPda()
	if err != nil {
		return nil, solana.PublicKey{}, err
	}
	var programConfig squads_multisig_program.ProgramConfig
	if err := client.GetAccountDataInto(ctx, programConfigPda, &programConfig); err != nil {
		return nil, solana.PublicKey{}, err
	}
	var configTreasury = programConfig.Treasury
	multisigPda, err := GetMultisigPda(createKey)
	if err != nil {
		return nil, solana.PublicKey{}, err
	}

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

// CreateMultisigTx creates a transaction for creating a multisig wallet
// Parameters:
// - ctx: Context for the operation
// - client: RPC client for Solana
// - createKey: One-time key used to generate multisigPda
// - creator: Signer and fee payer
// - configAuthority: Optional configuration authority
// - members: List of members for the multisig
// - threshold: Number of signatures required for approval
// - timelock: Time lock period
// - rentCollector: Optional rent collector
// Returns:
// - Transaction for creating a multisig wallet
// - Public key of the created multisig
// - Error, if any
func CreateMultisigTx(ctx context.Context, client *rpc.Client, createKey, creator solana.PublicKey, configAuthority *solana.PublicKey, members []squads_multisig_program.Member, threshold uint16, timelock uint32, rentCollector *solana.PublicKey) (*solana.Transaction, solana.PublicKey, error) {
	ix, multisigPda, err := CreateMultisigIx(ctx, client, createKey, creator, configAuthority, members, threshold, timelock, rentCollector)
	if err != nil {
		return nil, multisigPda, err
	}
	recent, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, multisigPda, err
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(creator),
	)
	return tx, multisigPda, err
}
