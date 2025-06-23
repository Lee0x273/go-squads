package squads

import (
	"context"

	"github.com/Lee0x273/go-squads/generated/squads_multisig_program"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
)

type Multisig struct {
	multisigPda solana.PublicKey
	client      *rpc.Client
}

func New(client *rpc.Client, multisigPda solana.PublicKey) *Multisig {
	return &Multisig{
		multisigPda: multisigPda,
		client:      client,
	}
}

func (s *Multisig) MultisigAccount(ctx context.Context) (*squads_multisig_program.Multisig, error) {
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

func (s *Multisig) VaultTransactionAccount(ctx context.Context, transactionPda solana.PublicKey) (*squads_multisig_program.VaultTransaction, error) {
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

func (s *Multisig) ProposalAccount(ctx context.Context, proposalPda solana.PublicKey) (*squads_multisig_program.Proposal, error) {
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

func (s *Multisig) VaultTransactionCreateIx(ctx context.Context, creatorAndPayer solana.PublicKey, vaultIndex uint8, transactionIndex uint64, instructions []solana.Instruction) (solana.Instruction, error) {
	vaultPda, err := GetVaultPda(s.multisigPda, vaultIndex)
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

	transactionPda, err := GetTransactionPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}

	txMessageBytes, err := TransactionMessageToMultisigTransactionMessageBytes(TransactionMessage{
		PayerKey:        vaultPda,
		Instructions:    instructions,
		RecentBlockhash: solana.Hash{}, //unused ,canbe zero hash
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

	return ix, nil
}

func (s *Multisig) VaultTransactionCreateTx(ctx context.Context, creatorAndPayer solana.PublicKey, vaultIndex uint8, transactionIndex uint64, instructions []solana.Instruction) (*solana.Transaction, error) {
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	ix, err := s.VaultTransactionCreateIx(ctx, creatorAndPayer, vaultIndex, transactionIndex, instructions)
	if err != nil {
		return nil, err
	}
	return solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(creatorAndPayer),
	)
}

func (s *Multisig) ProposalCreateIx(ctx context.Context, creatorAndPayer solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	args := squads_multisig_program.ProposalCreateArgs{
		TransactionIndex: transactionIndex,
		Draft:            false,
	}
	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	ix := squads_multisig_program.NewProposalCreateInstruction(
		args,
		s.multisigPda,
		proposalPda,
		creatorAndPayer, // creator
		creatorAndPayer,
		solana.SystemProgramID,
	).Build()
	return ix, nil
}

func (s *Multisig) ProposalCreateTx(ctx context.Context, creatorAndPayer solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	ix, err := s.ProposalCreateIx(ctx, creatorAndPayer, transactionIndex)
	if err != nil {
		return nil, err
	}

	return solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(creatorAndPayer),
	)
}

func (s *Multisig) VaultTransactionAndProposalTx(ctx context.Context, creatorAndPayer solana.PublicKey, vaultIndex uint8, transactionIndex uint64, instructions []solana.Instruction, autoApprove bool) (*solana.Transaction, error) {
	vaultPda, err := GetVaultPda(s.multisigPda, vaultIndex)
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

	transactionPda, err := GetTransactionPda(s.multisigPda, transactionIndex)
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

	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
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
		proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
		if err != nil {
			return nil, err
		}
		paIx := squads_multisig_program.NewProposalApproveInstruction(
			squads_multisig_program.ProposalVoteArgs{},
			s.multisigPda,
			creatorAndPayer,
			proposalPda,
		).Build()
		ixs = append(ixs, paIx)
	}

	return solana.NewTransaction(
		ixs,
		recent.Value.Blockhash,
		solana.TransactionPayer(creatorAndPayer),
	)
}

func (s *Multisig) ProposalApproveIx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	args := squads_multisig_program.ProposalVoteArgs{}

	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}

	ix := squads_multisig_program.NewProposalApproveInstruction(
		args,
		s.multisigPda,
		voter,
		proposalPda,
	).Build()

	return ix, nil
}

func (s *Multisig) ProposalApproveTx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	ix, err := s.ProposalApproveIx(ctx, voter, transactionIndex)
	if err != nil {
		return nil, err
	}
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	return solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(voter),
	)
}

func (s *Multisig) ProposalVoteIx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64, op VoteOP) (solana.Instruction, error) {
	args := squads_multisig_program.ProposalVoteArgs{}

	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	var ix solana.Instruction
	switch op {
	case VoteOPApprove:
		ix = squads_multisig_program.NewProposalApproveInstruction(
			args,
			s.multisigPda,
			voter,
			proposalPda,
		).Build()
	case VoteOPReject:
		ix = squads_multisig_program.NewProposalRejectInstruction(
			args,
			s.multisigPda,
			voter,
			proposalPda,
		).Build()
	case VoteOPCancel:
		// proposalAccount, err := s.ProposalAccount(ctx, proposalPda)
		// if err != nil {
		// 	return nil, err
		// }
		// if _, ok := proposalAccount.Status.(*squads_multisig_program.ProposalStatusApproved); !ok {
		// 	return nil, fmt.Errorf("proposal is not approved")
		// }
		ix = squads_multisig_program.NewProposalCancelInstruction(
			args,
			s.multisigPda,
			voter,
			proposalPda,
		).Build()
	}
	return ix, nil
}

func (s *Multisig) ProposalVoteTx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64, op VoteOP) (*solana.Transaction, error) {
	ix, err := s.ProposalVoteIx(ctx, voter, transactionIndex, op)
	if err != nil {
		return nil, err
	}
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	return solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(voter),
	)
}

func (s *Multisig) VaultTransactionExecuteIx(ctx context.Context, executor solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	transactionPda, err := GetTransactionPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}

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

	return ixb.Build(), nil
}

func (s *Multisig) VaultTransactionExecuteTx(ctx context.Context, executor solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	ix, err := s.VaultTransactionExecuteIx(ctx, executor, transactionIndex)
	if err != nil {
		return nil, err
	}
	recent, err := s.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	return solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(executor),
	)
}
