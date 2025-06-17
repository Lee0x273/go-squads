package squads

import (
	"context"
	"fmt"

	"github.com/Lee0x273/go-squads/generated/squads_multisig_program"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
)

type Squads struct {
	multisigPda solana.PublicKey
	client      *rpc.Client
}

func NewSquads(client *rpc.Client, multisigPda solana.PublicKey) *Squads {
	return &Squads{
		multisigPda: multisigPda,
		client:      client,
	}
}

func (s *Squads) MultisigAccount(ctx context.Context) (*squads_multisig_program.Multisig, error) {
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

func (s *Squads) VaultTransactionAccount(ctx context.Context, transactionPda solana.PublicKey) (*squads_multisig_program.VaultTransaction, error) {
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

func (s *Squads) ProposalAccount(ctx context.Context, proposalPda solana.PublicKey) (*squads_multisig_program.Proposal, error) {
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

func (s *Squads) CreateVaultTransactionCreateTx(ctx context.Context, creatorAndPayer solana.PublicKey, vaultIndex uint8, transactionIndex uint64, instructions []solana.Instruction) (*solana.Transaction, error) {
	vaultPda, err := s.GetVaultPda(vaultIndex)
	if err != nil {
		return nil, err
	}
	if transactionIndex == 0 {
		multisigInfo, err := s.MultisigAccount(ctx)
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

func (s *Squads) CreateProposalCreateTx(ctx context.Context, creatorAndPayer solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
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

func (s *Squads) CreateVaultTransactionAndProposalTx(ctx context.Context, creatorAndPayer solana.PublicKey, vaultIndex uint8, transactionIndex uint64, instructions []solana.Instruction, autoApprove bool) (*solana.Transaction, error) {
	vaultPda, err := s.GetVaultPda(vaultIndex)
	if err != nil {
		return nil, err
	}
	if transactionIndex == 0 {
		multisigInfo, err := s.MultisigAccount(ctx)
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

	vtcIx := squads_multisig_program.NewVaultTransactionCreateInstruction(
		args,
		s.multisigPda,
		transactionPda,
		creatorAndPayer, // creator
		creatorAndPayer,
		solana.SystemProgramID,
	).Build()

	proposalPda, _ := s.GetProposalPda(transactionIndex)
	pcIx := squads_multisig_program.NewProposalCreateInstruction(
		squads_multisig_program.ProposalCreateArgs{
			TransactionIndex: transactionIndex,
			Draft:            false,
		},
		s.multisigPda,
		proposalPda,
		creatorAndPayer,
		creatorAndPayer,
		solana.SystemProgramID,
	).Build()

	ixs := []solana.Instruction{vtcIx, pcIx}
	if autoApprove { // creator must be a voter
		proposalPda, _ := s.GetProposalPda(transactionIndex)
		paIx := squads_multisig_program.NewProposalApproveInstruction(
			squads_multisig_program.ProposalVoteArgs{},
			s.multisigPda,
			creatorAndPayer,
			proposalPda,
		).Build()
		ixs = append(ixs, paIx)
	}

	tx, _ := solana.NewTransaction(
		ixs,
		recent.Value.Blockhash,
		solana.TransactionPayer(creatorAndPayer),
	)
	return tx, nil
}

func (s *Squads) CreateProposalApproveTx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
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

func (s *Squads) CreateProposalVoteTx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64, op VoteOP) (*solana.Transaction, error) {
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	args := squads_multisig_program.ProposalVoteArgs{}

	proposalPda, _ := s.GetProposalPda(transactionIndex)

	ixs := []solana.Instruction{}

	switch op {
	case VoteOPApprove:
		ix := squads_multisig_program.NewProposalApproveInstruction(
			args,
			s.multisigPda,
			voter,
			proposalPda,
		).Build()
		ixs = append(ixs, ix)
	case VoteOPReject:
		ix := squads_multisig_program.NewProposalRejectInstruction(
			args,
			s.multisigPda,
			voter,
			proposalPda,
		).Build()
		ixs = append(ixs, ix)
	case VoteOPCancel:
		proposalAccount, err := s.ProposalAccount(ctx, proposalPda)
		if err != nil {
			return nil, err
		}
		if _, ok := proposalAccount.Status.(*squads_multisig_program.ProposalStatusApproved); !ok {
			return nil, fmt.Errorf("proposal is not approved")
		}
		ix := squads_multisig_program.NewProposalCancelInstruction(
			args,
			s.multisigPda,
			voter,
			proposalPda,
		).Build()
		ixs = append(ixs, ix)
	}

	tx, _ := solana.NewTransaction(
		ixs,
		recent.Value.Blockhash,
		solana.TransactionPayer(voter),
	)
	return tx, nil
}

func (s *Squads) CreateVaultTransactionExecuteTx(ctx context.Context, executor solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	transactionPda, err := s.GetTransactionPda(transactionIndex)
	if err != nil {
		return nil, err
	}
	proposalPda, _ := s.GetProposalPda(transactionIndex)

	vaultTransaction, err := s.VaultTransactionAccount(ctx, transactionPda)
	if err != nil {
		return nil, err
	}
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
