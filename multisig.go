package squads

import (
	"context"

	"github.com/Lee0x273/go-squads/generated/squads_multisig_program"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
)

// Multisig represents a multisig wallet
type Multisig struct {
	multisigPda solana.PublicKey
	client      *rpc.Client
}

// New creates a new Multisig instance
func New(client *rpc.Client, multisigPda solana.PublicKey) *Multisig {
	return &Multisig{
		multisigPda: multisigPda,
		client:      client,
	}
}

// MultisigAccount retrieves the multisig account information
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

// VaultTransactionAccount retrieves the vault transaction account information
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

// ProposalAccount retrieves the proposal account information
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

// MultisigAddMemeberIx creates an instruction to add a member to the multisig
func (s *Multisig) MultisigAddMemeberIx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, member squads_multisig_program.Member) (solana.Instruction, error) {
	args := squads_multisig_program.MultisigAddMemberArgs{
		NewMember: member,
	}

	ix := squads_multisig_program.NewMultisigAddMemberInstruction(
		args,
		s.multisigPda,
		configAuthority,
		rentPayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

// MultisigAddMemeberTx creates a transaction to add a member to the multisig
func (s *Multisig) MultisigAddMemeberTx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, member squads_multisig_program.Member) (*solana.Transaction, error) {
	ix, err := s.MultisigAddMemeberIx(ctx, configAuthority, rentPayer, member)
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
		solana.TransactionPayer(rentPayer),
	)
}

// MultisigRemoveMemberIx creates an instruction to remove a member from the multisig
func (s *Multisig) MultisigRemoveMemberIx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, member solana.PublicKey) (solana.Instruction, error) {
	args := squads_multisig_program.MultisigRemoveMemberArgs{
		OldMember: member,
	}

	ix := squads_multisig_program.NewMultisigRemoveMemberInstruction(
		args,
		s.multisigPda,
		configAuthority,
		rentPayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

// MultisigRemoveMemberTx creates a transaction to remove a member from the multisig
func (s *Multisig) MultisigRemoveMemberTx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, member solana.PublicKey) (*solana.Transaction, error) {
	ix, err := s.MultisigRemoveMemberIx(ctx, configAuthority, rentPayer, member)
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
		solana.TransactionPayer(rentPayer),
	)
}

// MultisigChangeThresholdIx creates an instruction to change the threshold of the multisig
func (s *Multisig) MultisigChangeThresholdIx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, threshold uint16) (solana.Instruction, error) {
	args := squads_multisig_program.MultisigChangeThresholdArgs{
		NewThreshold: threshold,
	}

	ix := squads_multisig_program.NewMultisigChangeThresholdInstruction(
		args,
		s.multisigPda,
		configAuthority,
		rentPayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

// MultisigChangeThresholdTx creates a transaction to change the threshold of the multisig
func (s *Multisig) MultisigChangeThresholdTx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, threshold uint16) (*solana.Transaction, error) {
	ix, err := s.MultisigChangeThresholdIx(ctx, configAuthority, rentPayer, threshold)
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
		solana.TransactionPayer(rentPayer),
	)
}

// MultisigSetConfigAuthorityIx creates an instruction to set the config authority of the multisig
func (s *Multisig) MultisigSetConfigAuthorityIx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, newConfigAuthority solana.PublicKey) (solana.Instruction, error) {
	args := squads_multisig_program.MultisigSetConfigAuthorityArgs{
		ConfigAuthority: newConfigAuthority,
	}

	ix := squads_multisig_program.NewMultisigSetConfigAuthorityInstruction(
		args,
		s.multisigPda,
		configAuthority,
		rentPayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

// MultisigSetConfigAuthorityTx creates a transaction to set the config authority of the multisig
func (s *Multisig) MultisigSetConfigAuthorityTx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, newConfigAuthority solana.PublicKey) (*solana.Transaction, error) {
	ix, err := s.MultisigSetConfigAuthorityIx(ctx, configAuthority, rentPayer, newConfigAuthority)
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
		solana.TransactionPayer(rentPayer),
	)
}

// MultisigSetRentCollectorIx creates an instruction to set the rent collector of the multisig
func (s *Multisig) MultisigSetRentCollectorIx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, rentCollector *solana.PublicKey) (solana.Instruction, error) {
	args := squads_multisig_program.MultisigSetRentCollectorArgs{
		RentCollector: rentCollector,
	}

	ix := squads_multisig_program.NewMultisigSetRentCollectorInstruction(
		args,
		s.multisigPda,
		configAuthority,
		rentPayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

// MultisigSetRentCollectorTx creates a transaction to set the rent collector of the multisig
func (s *Multisig) MultisigSetRentCollectorTx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, rentCollector *solana.PublicKey) (*solana.Transaction, error) {
	ix, err := s.MultisigSetRentCollectorIx(ctx, configAuthority, rentPayer, rentCollector)
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
		solana.TransactionPayer(rentPayer),
	)
}

// MultisigSetTimeLockIx creates an instruction to set the time lock for the multisig
func (s *Multisig) MultisigSetTimeLockIx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, timelock uint32) (solana.Instruction, error) {
	args := squads_multisig_program.MultisigSetTimeLockArgs{
		TimeLock: timelock,
	}

	ix := squads_multisig_program.NewMultisigSetTimeLockInstruction(
		args,
		s.multisigPda,
		configAuthority,
		rentPayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

// MultisigSetTimeLockTx creates a transaction to set the time lock for the multisig
func (s *Multisig) MultisigSetTimeLockTx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, timelock uint32) (*solana.Transaction, error) {
	ix, err := s.MultisigSetTimeLockIx(ctx, configAuthority, rentPayer, timelock)
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
		solana.TransactionPayer(rentPayer),
	)
}

// MultisigAddSpendingLimitIx creates an instruction to add a spending limit to the multisig
func (s *Multisig) MultisigAddSpendingLimitIx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, spendingLimitPda solana.PublicKey, args *squads_multisig_program.MultisigAddSpendingLimitArgs) (solana.Instruction, error) {
	ix := squads_multisig_program.NewMultisigAddSpendingLimitInstruction(
		*args,
		s.multisigPda,
		configAuthority,
		spendingLimitPda,
		rentPayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

// MultisigAddSpendingLimitTx creates a transaction to add a spending limit to the multisig
func (s *Multisig) MultisigAddSpendingLimitTx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, spendingLimitPda solana.PublicKey, args *squads_multisig_program.MultisigAddSpendingLimitArgs) (*solana.Transaction, error) {
	ix, err := s.MultisigAddSpendingLimitIx(ctx, configAuthority, rentPayer, spendingLimitPda, args)
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
		solana.TransactionPayer(rentPayer),
	)
}

// MultisigRemoveSpendingLimitIx creates an instruction to remove a spending limit from the multisig
func (s *Multisig) MultisigRemoveSpendingLimitIx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, spendingLimitPda solana.PublicKey) (solana.Instruction, error) {
	args := squads_multisig_program.MultisigRemoveSpendingLimitArgs{}
	ix := squads_multisig_program.NewMultisigRemoveSpendingLimitInstruction(
		args,
		s.multisigPda,
		configAuthority,
		spendingLimitPda,
		rentPayer,
	).Build()

	return ix, nil
}

// MultisigRemoveSpendingLimitTx creates a transaction to remove a spending limit from the multisig
func (s *Multisig) MultisigRemoveSpendingLimitTx(ctx context.Context, configAuthority, rentPayer solana.PublicKey, spendingLimitPda solana.PublicKey) (*solana.Transaction, error) {
	ix, err := s.MultisigRemoveSpendingLimitIx(ctx, configAuthority, rentPayer, spendingLimitPda)
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
		solana.TransactionPayer(rentPayer),
	)
}

// SpendingLimitUseIx creates an instruction to use a spending limit from a multisig vault
func (s *Multisig) SpendingLimitUseIx(ctx context.Context, member, spendingLimitPda, vault, destination, mint, vaultTokenAccount, destinationTokenAccount solana.PublicKey, args *squads_multisig_program.SpendingLimitUseArgs) (solana.Instruction, error) {
	ix := squads_multisig_program.NewSpendingLimitUseInstruction(
		*args,
		s.multisigPda,
		member,
		spendingLimitPda,
		vault,
		destination,
		solana.SystemProgramID,
		mint,
		vaultTokenAccount,
		destinationTokenAccount,
		solana.TokenProgramID,
	).Build()

	return ix, nil
}

// SpendingLimitUseTx creates a complete transaction to use a spending limit from a multisig vault
func (s *Multisig) SpendingLimitUseTx(ctx context.Context, member, spendingLimitPda, vault, destination, mint, vaultTokenAccount, destinationTokenAccount solana.PublicKey, args *squads_multisig_program.SpendingLimitUseArgs) (*solana.Transaction, error) {
	ix, err := s.SpendingLimitUseIx(ctx, member, spendingLimitPda, vault, destination, mint, vaultTokenAccount, destinationTokenAccount, args)
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
		solana.TransactionPayer(member),
	)
}

// VaultTransactionCreateIx creates an instruction to create a vault transaction
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

// VaultTransactionCreateTx creates a transaction to create a vault transaction
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

// ProposalCreateIx creates an instruction to create a proposal
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

// ProposalCreateTx creates a transaction to create a proposal
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

// VaultTransactionAndProposalTx creates a transaction that includes both vault transaction creation and proposal creation
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
	if err != nil {
		return nil, err
	}

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

// ProposalApproveIx creates an instruction to approve a proposal
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

// ProposalApproveTx creates a transaction to approve a proposal
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

// ProposalRejectIx creates an instruction to reject a proposal.
func (s *Multisig) ProposalRejectIx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	args := squads_multisig_program.ProposalVoteArgs{}

	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}

	ix := squads_multisig_program.NewProposalRejectInstruction(
		args,
		s.multisigPda,
		voter,
		proposalPda,
	).Build()

	return ix, nil
}

// ProposalRejectTx creates a transaction to reject a proposal.
func (s *Multisig) ProposalRejectTx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	ix, err := s.ProposalRejectIx(ctx, voter, transactionIndex)
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

// ProposalCancelIx creates an instruction to cancel a proposal.
func (s *Multisig) ProposalCancelIx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	args := squads_multisig_program.ProposalVoteArgs{}

	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}

	ix := squads_multisig_program.NewProposalCancelInstruction(
		args,
		s.multisigPda,
		voter,
		proposalPda,
	).Build()

	return ix, nil
}

// ProposalCancelTx creates a transaction to cancel a proposal.
func (s *Multisig) ProposalCancelTx(ctx context.Context, voter solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	ix, err := s.ProposalCancelIx(ctx, voter, transactionIndex)
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

// VaultTransactionExecuteIx creates an instruction to execute a vault transaction
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

// VaultTransactionExecuteTx creates a transaction to execute a vault transaction
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

// VaultTransactionAccountsCloseIx creates an instruction to close transaction accounts associated with a vault
func (s *Multisig) VaultTransactionAccountsCloseIx(ctx context.Context, feePayer solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}

	transactionPda, err := GetTransactionPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}

	ix := squads_multisig_program.NewVaultTransactionAccountsCloseInstruction(
		s.multisigPda,
		proposalPda,
		transactionPda,
		feePayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

// VaultTransactionAccountsCloseTx creates a complete transaction to close transaction accounts associated with a vault
func (s *Multisig) VaultTransactionAccountsCloseTx(ctx context.Context, feePayer solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	ix, err := s.VaultTransactionAccountsCloseIx(ctx, feePayer, transactionIndex)
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
		solana.TransactionPayer(feePayer),
	)
}

func (s *Multisig) ConfigTransactionCreateIx(ctx context.Context, creator solana.PublicKey, transactionIndex uint64, args *squads_multisig_program.ConfigTransactionCreateArgs) (solana.Instruction, error) {
	transactionPda, err := GetTransactionPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	ix := squads_multisig_program.NewConfigTransactionCreateInstruction(
		*args,
		s.multisigPda,
		transactionPda,
		creator,
		creator,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

func (s *Multisig) ConfigTransactionCreateTx(ctx context.Context, creator solana.PublicKey, transactionIndex uint64, args *squads_multisig_program.ConfigTransactionCreateArgs) (*solana.Transaction, error) {
	ix, err := s.ConfigTransactionCreateIx(ctx, creator, transactionIndex, args)
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
		solana.TransactionPayer(creator),
	)
}

func (s *Multisig) ConfigTransactionExecuteIx(ctx context.Context, member, feePayer solana.PublicKey, transactionIndex uint64, args *squads_multisig_program.ConfigTransactionCreateArgs) (solana.Instruction, error) {
	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	transactionPda, err := GetTransactionPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	ix := squads_multisig_program.NewConfigTransactionExecuteInstruction(
		s.multisigPda,
		member,
		proposalPda,
		transactionPda,
		feePayer,
		solana.SystemProgramID,
	).Build()

	return ix, nil
}

func (s *Multisig) ConfigTransactionExecuteTx(ctx context.Context, member, feePayer solana.PublicKey, transactionIndex uint64, args *squads_multisig_program.ConfigTransactionCreateArgs) (*solana.Transaction, error) {
	ix, err := s.ConfigTransactionExecuteIx(ctx, member, feePayer, transactionIndex, args)
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
		solana.TransactionPayer(feePayer),
	)
}

// ProposalActivateIx creates an instruction to activate a proposal.
func (s *Multisig) ProposalActivateIx(ctx context.Context, member solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	ix := squads_multisig_program.NewProposalActivateInstruction(
		s.multisigPda,
		member,
		proposalPda,
	).Build()
	return ix, nil
}

// ProposalActivateTx creates a transaction to activate a proposal.
func (s *Multisig) ProposalActivateTx(ctx context.Context, member solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	ix, err := s.ProposalActivateIx(ctx, member, transactionIndex)
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
		solana.TransactionPayer(member),
	)
}

// ProposalCancelV2Ix creates an instruction to cancel a proposal using the V2 instruction.
func (s *Multisig) ProposalCancelV2Ix(ctx context.Context, member solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	ix := squads_multisig_program.NewProposalCancelV2Instruction(
		squads_multisig_program.ProposalVoteArgs{},
		s.multisigPda,
		member,
		proposalPda,
		solana.SystemProgramID,
	).Build()
	return ix, nil
}

// ProposalCancelV2Tx creates a transaction to cancel a proposal using the V2 instruction.
func (s *Multisig) ProposalCancelV2Tx(ctx context.Context, member solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	ix, err := s.ProposalCancelV2Ix(ctx, member, transactionIndex)
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
		solana.TransactionPayer(member),
	)
}

// ConfigTransactionAccountsCloseIx creates an instruction to close a config transaction and its proposal.
func (s *Multisig) ConfigTransactionAccountsCloseIx(ctx context.Context, rentCollector solana.PublicKey, transactionIndex uint64) (solana.Instruction, error) {
	proposalPda, err := GetProposalPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	transactionPda, err := GetTransactionPda(s.multisigPda, transactionIndex)
	if err != nil {
		return nil, err
	}
	ix := squads_multisig_program.NewConfigTransactionAccountsCloseInstruction(
		s.multisigPda,
		proposalPda,
		transactionPda,
		rentCollector,
		solana.SystemProgramID,
	).Build()
	return ix, nil
}

// ConfigTransactionAccountsCloseTx creates a transaction to close a config transaction and its proposal.
func (s *Multisig) ConfigTransactionAccountsCloseTx(ctx context.Context, rentCollector solana.PublicKey, transactionIndex uint64) (*solana.Transaction, error) {
	ix, err := s.ConfigTransactionAccountsCloseIx(ctx, rentCollector, transactionIndex)
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
		solana.TransactionPayer(rentCollector),
	)
}
