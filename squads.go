package squads

import (
	"context"
	"fmt"
	"squads/generated/squads_multisig_program"

	"github.com/axengine/utils"
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

	account := &squads_multisig_program.Multisig{}
	decoder := ag_binary.NewBorshDecoder(data)
	if err := account.UnmarshalWithDecoder(decoder); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *SQuard) ProposalAccount(ctx context.Context, proposalPda solana.PublicKey) (*squads_multisig_program.Proposal, error) {
	out, err := s.client.GetAccountInfo(ctx, proposalPda)
	if err != nil {
		return nil, err
	}
	data := out.Value.Data.GetBinary()

	account := &squads_multisig_program.Proposal{}
	decoder := ag_binary.NewBorshDecoder(data)
	if err := account.UnmarshalWithDecoder(decoder); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *SQuard) VaultTransactionAccount(ctx context.Context, transactionPda solana.PublicKey) (*squads_multisig_program.VaultTransaction, error) {
	out, err := s.client.GetAccountInfo(ctx, transactionPda)
	if err != nil {
		return nil, err
	}
	data := out.Value.Data.GetBinary()

	account := &squads_multisig_program.VaultTransaction{}
	decoder := ag_binary.NewBorshDecoder(data)
	if err := account.UnmarshalWithDecoder(decoder); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *SQuard) CreateVaultTransactionCreate(ctx context.Context, creatorAndPayer solana.PublicKey, vaultIndex uint8, transactionIndex uint64, instructions []solana.Instruction) (*solana.Transaction, error) {
	vaultPda, err := s.GetVaultPda(vaultIndex)
	if err != nil {
		return nil, err
	}
	fmt.Println("vaultPda:", vaultPda)
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
	fmt.Println("transactionPda:", transactionPda)
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	txMessageBytes, err := TransactionMessageToMultisigTransactionMessageBytes(TransactionMessage{
		PayerKey:        vaultPda,
		Instructions:    instructions,
		RecentBlockhash: recent.Value.Blockhash,
	}, []addresslookuptable.KeyedAddressLookupTable{})
	fmt.Println("txMessageBytes len:", len(txMessageBytes), txMessageBytes)

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
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	args := squads_multisig_program.ProposalCreateArgs{
		TransactionIndex: transactionIndex,
		Draft:            false,
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
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	args := squads_multisig_program.ProposalVoteArgs{}

	proposalPda, _ := s.GetProposalPda(transactionIndex)

	ix := squads_multisig_program.NewProposalApproveInstruction(
		args,
		s.multisigPda,
		voter,
		proposalPda,
	).Build()

	tx, _ := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(voter),
	)
	return tx, nil
}

func (s *SQuard) CreateVaultTransactionExecute(ctx context.Context, executor solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	transactionPda, err := s.GetTransactionPda(transactionIndex)
	if err != nil {
		return nil, err
	}
	proposalPda, _ := s.GetProposalPda(transactionIndex)

	vaultTransaction, err := s.VaultTransactionAccount(ctx, transactionPda)
	if err != nil {
		return nil, err
	}
	utils.JsonPrettyToStdout(vaultTransaction.Message)
	additionalAccounts := vaultTransaction.Message.AccountKeys

	ixb := squads_multisig_program.NewVaultTransactionExecuteInstruction(
		s.multisigPda,
		proposalPda,
		transactionPda,
		executor,
	)

	// Append additional accounts with dynamic properties
	for i, accountKey := range additionalAccounts {
		isWritable := false
		// Determine if the account is writable based on the message structure
		if i < int(vaultTransaction.Message.NumWritableSigners) {
			isWritable = true // Writable signer
		} else if (i - int(vaultTransaction.Message.NumSigners)) < int(vaultTransaction.Message.NumWritableNonSigners) {
			isWritable = true // Writable non-signer
		}
		// Additional accounts are typically not signers in multisig execution
		ixb.AccountMetaSlice = append(ixb.AccountMetaSlice,
			solana.NewAccountMeta(accountKey, isWritable, false))
	}
	ix := ixb.Build()

	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	tx, _ := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(executor),
	)
	return tx, nil
}
