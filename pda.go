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

func (s *SQuard) GetProgramConfigPda() (solana.PublicKey, error) {
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

func (s *SQuard) GetMultisigPda(createKey solana.PublicKey) (solana.PublicKey, error) {
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

func (s *SQuard) GetVaultPda(index uint8) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			s.multisigPda.Bytes(),
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

func (s *SQuard) GetEphemeralSignerPda(transactionPda solana.PublicKey, ephemeralSignerIndex uint8) (solana.PublicKey, error) {
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

func (s *SQuard) GetTransactionPda(index uint64) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			s.multisigPda.Bytes(),
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

func (s *SQuard) GetProposalPda(transactionIndex uint64) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			s.multisigPda.Bytes(),
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

func (s *SQuard) GetBatchTransactionPda(batchIndex uint64, transactionIndex uint32) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			s.multisigPda.Bytes(),
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

func (s *SQuard) GetSpendingLimitPda(createKey solana.PublicKey) (solana.PublicKey, error) {
	pk, _, err := solana.FindProgramAddress(
		[][]byte{
			SEED_PREFIX,
			s.multisigPda.Bytes(),
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
