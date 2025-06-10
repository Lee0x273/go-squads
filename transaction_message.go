package squads

import (
	"bytes"
	"log"
	"squads/generated/squads_multisig_program"

	ag_binary "github.com/gagliardetto/binary"
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

	txMsg := squads_multisig_program.VaultTransactionMessage{
		NumSigners:            uint8(compiledMessage.Header.NumRequiredSignatures),
		NumWritableSigners:    uint8(compiledMessage.Header.NumRequiredSignatures - compiledMessage.Header.NumReadonlySignedAccounts),
		NumWritableNonSigners: uint8(len(compiledMessage.StaticAccountKeys)) - compiledMessage.Header.NumRequiredSignatures - compiledMessage.Header.NumReadonlyUnsignedAccounts,
		AccountKeys:           compiledMessage.StaticAccountKeys,
		// Instructions:          compiledMessage.Instructions,
		// AddressTableLookups:   compiledMessage.AddressTableLookups,
	}
	for _, v := range compiledMessage.CompiledInstructions {
		txMsg.Instructions = append(txMsg.Instructions, squads_multisig_program.MultisigCompiledInstruction{
			ProgramIdIndex: uint8(v.ProgramIDIndex),
			AccountIndexes: convertToUint8Slice(v.Accounts),
			Data:           v.Data,
		})
	}

	// Serialize the message to bytes using borsh encoder
	buf := new(bytes.Buffer)
	err := ag_binary.NewBorshEncoder(buf).Encode(txMsg)
	if err != nil {
		log.Fatalf("Failed to encode transaction message: %v", err)
	}

	return buf.Bytes(), nil
}
