package squads

import (
	"bytes"
	"squads/generated/squads_multisig_program"

	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
)

// TransactionMessage represents the transaction message structure
type TransactionMessage struct {
	PayerKey        solana.PublicKey
	Instructions    []solana.Instruction
	RecentBlockhash solana.Hash
}

// TransactionMessageToMultisigTransactionMessageBytes converts a transaction message to bytes
func TransactionMessageToMultisigTransactionMessageBytes(message TransactionMessage,
	addressLookupTableAccounts []addresslookuptable.KeyedAddressLookupTable) ([]byte, error) {
	// Compile the message to V0 format
	compiledMessage := CompileToWrappedMessageV0(message.PayerKey,
		message.RecentBlockhash,
		message.Instructions,
		addressLookupTableAccounts)
	txMsg := squads_multisig_program.TransactionMessage{
		NumSigners:            uint8(compiledMessage.Header.NumRequiredSignatures),
		NumWritableSigners:    uint8(compiledMessage.Header.NumRequiredSignatures - compiledMessage.Header.NumReadonlySignedAccounts),
		NumWritableNonSigners: uint8(len(compiledMessage.AccountKeys)) - compiledMessage.Header.NumRequiredSignatures - compiledMessage.Header.NumReadonlyUnsignedAccounts,
		AccountKeys: squads_multisig_program.SmallVec[uint8, solana.PublicKey]{
			Data: compiledMessage.AccountKeys,
		},
		Instructions:        squads_multisig_program.SmallVec[uint8, squads_multisig_program.CompiledInstruction]{},
		AddressTableLookups: squads_multisig_program.SmallVec[uint8, squads_multisig_program.MessageAddressTableLookup]{},
	}
	for _, v := range compiledMessage.Instructions {
		txMsg.Instructions.Data = append(txMsg.Instructions.Data, squads_multisig_program.CompiledInstruction{
			ProgramIdIndex: uint8(v.ProgramIDIndex),
			AccountIndexes: squads_multisig_program.SmallVec[uint8, uint8]{Data: convertToUint8Slice(v.Accounts)},
			Data:           squads_multisig_program.SmallVec[uint16, uint8]{Data: v.Data},
		})
	}
	for _, v := range compiledMessage.AddressTableLookups {
		txMsg.AddressTableLookups.Data = append(txMsg.AddressTableLookups.Data, squads_multisig_program.MessageAddressTableLookup{
			AccountKey:      v.AccountKey,
			WritableIndexes: squads_multisig_program.SmallVec[uint8, uint8]{Data: v.WritableIndexes},
			ReadonlyIndexes: squads_multisig_program.SmallVec[uint8, uint8]{Data: v.ReadonlyIndexes},
		})
	}

	// encode custom
	buf := new(bytes.Buffer)
	if err := squads_multisig_program.NewEncoder(buf).Encode(&txMsg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
