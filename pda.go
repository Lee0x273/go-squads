package squads

import (
	"github.com/Lee0x273/go-squads/generated/squads_multisig_program"
	"github.com/gagliardetto/solana-go"
)

var (
	SEED_PREFIX            = []byte("multisig")
	SEED_PROGRAM_CONFIG    = []byte("program_config")
	SEED_MULTISIG          = []byte("multisig")
	SEED_VAULT             = []byte("vault")
	SEED_TRANSACTION       = []byte("transaction")
	SEED_PROPOSAL          = []byte("proposal")
	SEED_BATCH_TRANSACTION = []byte("batch_transaction")
	SEED_EPHEMERAL_SIGNER  = []byte("ephemeral_signer")
	SEED_SPENDING_LIMIT    = []byte("spending_limit")
)

func GetProgramConfigPda() (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			SEED_PROGRAM_CONFIG,
		},
		squads_multisig_program.ProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return pk, nil
}

func GetMultisigPda(createKey solana.PublicKey) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			SEED_MULTISIG,
			createKey.Bytes(),
		},
		squads_multisig_program.ProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return pk, nil
}

func GetVaultPda(multisigPda solana.PublicKey, index uint8) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			multisigPda.Bytes(),
			SEED_VAULT,
			[]byte{index},
		},
		squads_multisig_program.ProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return pk, nil
}

func GetEphemeralSignerPda(transactionPda solana.PublicKey, ephemeralSignerIndex uint8) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			transactionPda.Bytes(),
			SEED_EPHEMERAL_SIGNER,
			[]byte{ephemeralSignerIndex},
		},
		squads_multisig_program.ProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return pk, nil
}

func GetTransactionPda(multisigPda solana.PublicKey, index uint64) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			multisigPda.Bytes(),
			SEED_TRANSACTION,
			toU64Bytes(index),
		},
		squads_multisig_program.ProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return pk, nil
}

func GetProposalPda(multisigPda solana.PublicKey, transactionIndex uint64) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			multisigPda.Bytes(),
			SEED_TRANSACTION,
			toU64Bytes(transactionIndex),
			SEED_PROPOSAL,
		},
		squads_multisig_program.ProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return pk, nil
}

func GetBatchTransactionPda(multisigPda solana.PublicKey, batchIndex uint64, transactionIndex uint32) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			multisigPda.Bytes(),
			SEED_TRANSACTION,
			toU64Bytes(batchIndex),
			SEED_BATCH_TRANSACTION,
			toU32Bytes(transactionIndex),
		},
		squads_multisig_program.ProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return pk, nil
}

func GetSpendingLimitPda(multisigPda solana.PublicKey, createKey solana.PublicKey) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			multisigPda.Bytes(),
			SEED_SPENDING_LIMIT,
			createKey.Bytes(),
		},
		squads_multisig_program.ProgramID,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}
	return pk, nil
}
