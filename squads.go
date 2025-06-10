package squads

import (
	"context"
	"squads/generated/squads_multisig_program"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
)

type SQuard struct {
	multisigPda solana.PublicKey
	client      *rpc.Client
}

func NewSQuard(multisigPda solana.PublicKey, client *rpc.Client) *SQuard {
	return &SQuard{
		multisigPda: multisigPda,
		client:      client,
	}
}

func (s *SQuard) Multisig(ctx context.Context) (*squads_multisig_program.Multisig, error) {
	out, err := s.client.GetAccountInfo(ctx, s.multisigPda)
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

func (s *SQuard) ProposalAccount(ctx context.Context, proposalPda solana.PublicKey) (*squads_multisig_program.Proposal, error) {
	out, err := s.client.GetAccountInfo(ctx, proposalPda)
	if err != nil {
		return nil, err
	}
	data := out.Value.Data.GetBinary()

	multisig := &squads_multisig_program.Proposal{}
	decoder := ag_binary.NewBorshDecoder(data)
	if err := multisig.UnmarshalWithDecoder(decoder); err != nil {
		return nil, err
	}
	return multisig, nil
}

func (s *SQuard) CreateVaultTransactionCreate(ctx context.Context, creatorAndPayer solana.PublicKey, vaultIndex uint8, transactionIndex uint64, instructions []solana.Instruction) (*solana.Transaction, error) {
	vaultPda, err := s.GetVaultPda(vaultIndex)
	if err != nil {
		return nil, err
	}
	if transactionIndex == 0 {
		multisigInfo, err := s.Multisig(ctx)
		if err != nil {
			return nil, err
		}
		transactionIndex = multisigInfo.TransactionIndex + 1
	}

	transactionPda, err := s.GetTransactionPda(transactionIndex)
	if err != nil {
		return nil, err
	}
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	txMessageBytes, err := TransactionMessageToMultisigTransactionMessageBytes(TransactionMessage{
		PayerKey:        vaultPda,
		Instructions:    instructions,
		RecentBlockhash: recent.Value.Blockhash,
	}, []addresslookuptable.KeyedAddressLookupTable{})

	args := squads_multisig_program.VaultTransactionCreateArgs{
		VaultIndex:         vaultIndex,
		EphemeralSigners:   0,
		TransactionMessage: txMessageBytes,
	}

	ix := squads_multisig_program.NewVaultTransactionCreateInstruction(
		args,
		s.multisigPda,
		transactionPda,
		creatorAndPayer, // creator
		creatorAndPayer,
		solana.SystemProgramID,
	).Build()

	tx, _ := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(creatorAndPayer),
	)
	return tx, nil
}

func (s *SQuard) CreateProposalCreate(ctx context.Context, creatorAndPayer solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}
	args := squads_multisig_program.ProposalCreateArgs{
		TransactionIndex: transactionIndex,
	}

	proposalPda, _ := s.GetProposalPda(transactionIndex)

	ix := squads_multisig_program.NewProposalCreateInstruction(
		args,
		s.multisigPda,
		proposalPda,
		creatorAndPayer, // creator
		creatorAndPayer,
		solana.SystemProgramID,
	).Build()

	tx, _ := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(creatorAndPayer),
	)
	return tx, nil
}

func (s *SQuard) CreateProposalVote(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}
	args := squads_multisig_program.ProposalVoteArgs{}

	proposalPda, _ := s.GetProposalPda(transactionIndex)

	ix := squads_multisig_program.NewProposalCancelV2Instruction(
		args,
		s.multisigPda,
		voter,
		proposalPda,
		solana.SystemProgramID,
	).Build()

	tx, _ := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(voter),
	)
	return tx, nil
}

func (s *SQuard) CreateVaultTransactionExecute(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}

	transactionPda, err := s.GetTransactionPda(transactionIndex)
	if err != nil {
		return nil, err
	}
	proposalPda, _ := s.GetProposalPda(transactionIndex)

	ix := squads_multisig_program.NewVaultTransactionExecuteInstruction(
		s.multisigPda,
		proposalPda,
		transactionPda,
		voter,
	).Build()

	tx, _ := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(voter),
	)
	return tx, nil
}
