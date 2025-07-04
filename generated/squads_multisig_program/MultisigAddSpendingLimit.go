// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package squads_multisig_program

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Create a new spending limit for the controlled multisig.
type MultisigAddSpendingLimit struct {
	Args *MultisigAddSpendingLimitArgs

	// [0] = [] multisig
	//
	// [1] = [SIGNER] configAuthority
	// ··········· Multisig `config_authority` that must authorize the configuration change.
	//
	// [2] = [WRITE] spendingLimit
	//
	// [3] = [WRITE, SIGNER] rentPayer
	// ··········· This is usually the same as `config_authority`, but can be a different account if needed.
	//
	// [4] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewMultisigAddSpendingLimitInstructionBuilder creates a new `MultisigAddSpendingLimit` instruction builder.
func NewMultisigAddSpendingLimitInstructionBuilder() *MultisigAddSpendingLimit {
	nd := &MultisigAddSpendingLimit{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 5),
	}
	return nd
}

// SetArgs sets the "args" parameter.
func (inst *MultisigAddSpendingLimit) SetArgs(args MultisigAddSpendingLimitArgs) *MultisigAddSpendingLimit {
	inst.Args = &args
	return inst
}

// SetMultisigAccount sets the "multisig" account.
func (inst *MultisigAddSpendingLimit) SetMultisigAccount(multisig ag_solanago.PublicKey) *MultisigAddSpendingLimit {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(multisig)
	return inst
}

// GetMultisigAccount gets the "multisig" account.
func (inst *MultisigAddSpendingLimit) GetMultisigAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetConfigAuthorityAccount sets the "configAuthority" account.
// Multisig `config_authority` that must authorize the configuration change.
func (inst *MultisigAddSpendingLimit) SetConfigAuthorityAccount(configAuthority ag_solanago.PublicKey) *MultisigAddSpendingLimit {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(configAuthority).SIGNER()
	return inst
}

// GetConfigAuthorityAccount gets the "configAuthority" account.
// Multisig `config_authority` that must authorize the configuration change.
func (inst *MultisigAddSpendingLimit) GetConfigAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSpendingLimitAccount sets the "spendingLimit" account.
func (inst *MultisigAddSpendingLimit) SetSpendingLimitAccount(spendingLimit ag_solanago.PublicKey) *MultisigAddSpendingLimit {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(spendingLimit).WRITE()
	return inst
}

// GetSpendingLimitAccount gets the "spendingLimit" account.
func (inst *MultisigAddSpendingLimit) GetSpendingLimitAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetRentPayerAccount sets the "rentPayer" account.
// This is usually the same as `config_authority`, but can be a different account if needed.
func (inst *MultisigAddSpendingLimit) SetRentPayerAccount(rentPayer ag_solanago.PublicKey) *MultisigAddSpendingLimit {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(rentPayer).WRITE().SIGNER()
	return inst
}

// GetRentPayerAccount gets the "rentPayer" account.
// This is usually the same as `config_authority`, but can be a different account if needed.
func (inst *MultisigAddSpendingLimit) GetRentPayerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *MultisigAddSpendingLimit) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *MultisigAddSpendingLimit {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *MultisigAddSpendingLimit) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

func (inst MultisigAddSpendingLimit) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_MultisigAddSpendingLimit,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst MultisigAddSpendingLimit) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *MultisigAddSpendingLimit) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Args == nil {
			return errors.New("Args parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Multisig is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.ConfigAuthority is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.SpendingLimit is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.RentPayer is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *MultisigAddSpendingLimit) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("MultisigAddSpendingLimit")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Args", *inst.Args))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=5]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("       multisig", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("configAuthority", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("  spendingLimit", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("      rentPayer", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("  systemProgram", inst.AccountMetaSlice.Get(4)))
					})
				})
		})
}

func (obj MultisigAddSpendingLimit) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Args` param:
	err = encoder.Encode(obj.Args)
	if err != nil {
		return err
	}
	return nil
}
func (obj *MultisigAddSpendingLimit) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Args`:
	err = decoder.Decode(&obj.Args)
	if err != nil {
		return err
	}
	return nil
}

// NewMultisigAddSpendingLimitInstruction declares a new MultisigAddSpendingLimit instruction with the provided parameters and accounts.
func NewMultisigAddSpendingLimitInstruction(
	// Parameters:
	args MultisigAddSpendingLimitArgs,
	// Accounts:
	multisig ag_solanago.PublicKey,
	configAuthority ag_solanago.PublicKey,
	spendingLimit ag_solanago.PublicKey,
	rentPayer ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *MultisigAddSpendingLimit {
	return NewMultisigAddSpendingLimitInstructionBuilder().
		SetArgs(args).
		SetMultisigAccount(multisig).
		SetConfigAuthorityAccount(configAuthority).
		SetSpendingLimitAccount(spendingLimit).
		SetRentPayerAccount(rentPayer).
		SetSystemProgramAccount(systemProgram)
}
