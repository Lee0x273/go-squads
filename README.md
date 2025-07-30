# Squads V4 Golang SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/Lee0x273/go-squads.svg)](https://pkg.go.dev/github.com/Lee0x273/go-squads)
[![Go Report Card](https://goreportcard.com/badge/github.com/Lee0x273/go-squads)](https://goreportcard.com/report/github.com/Lee0x273/go-squads)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This is an unofficial Golang SDK for the [Squads Protocol V4](https://squads.so/). It provides a convenient way to interact with the Squads multisig program on Solana.

## Installation

```bash
go get github.com/Lee0x273/go-squads
```

## Usage

### Initialize the Client

First, you need to create a new Squads client with an RPC connection and the public key of your multisig account.

```go
import (
	"context"
	"fmt"
	"github.com/Lee0x273/go-squads"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	client := rpc.New(rpc.DevNet.RPC)
	multisigPda := solana.MustPublicKeyFromBase58("YOUR_MULTISIG_PDA")
	s := squads.New(client, multisigPda)

	// You can now use the 's' object to interact with your multisig.
}
```

### Create a Multisig

You can create a new multisig account with a set of members and a threshold.

```go
// ... (inside a function with your RPC client)

// The creator of the multisig.
creator := // Your creator's public key

// A new, unique keypair for the multisig.
createKey := solana.NewWallet()

// The members of the multisig, with their permissions.
// In this example, both members have full permissions.
members := []squads_multisig_program.Member{
    {Key: member1, Permissions: squads.Permissions{Mask: uint8(squads.Initiate | squads.Vote | squads.Execute)}},
    {Key: member2, Permissions: squads.Permissions{Mask: uint8(squads.Initiate | squads.Vote | squads.Execute)}},
}

// The number of signatures required for a transaction to be approved.
threshold := uint16(2)

// The time lock period in seconds.
timelock := uint32(0)

// Create the multisig transaction.
tx, multisigPda, err := squads.CreateMultisigTx(
    context.Background(),
    client,
    createKey.PublicKey(),
    creator,
    nil, // configAuthority
    members,
    threshold,
    timelock,
    nil, // rentCollector
)
if err != nil {
    // Handle error
}

// Sign and send the transaction...
fmt.Printf("Multisig created: %s\n", multisigPda.String())
```


### Create a Vault Transaction

You can create and sign a transaction to interact with a vault.

```go
// ... (inside a function with your signer and squads client)

// The public key of the account signing the transaction.
signer := // Your signer's public key

// The vault index you want to interact with.
vaultIndex := uint64(0)

// The instruction to execute on the vault.
// For example, a simple transfer instruction.
vaultInstruction := system.NewTransferInstruction(
    100000000, // Amount in lamports
    vaultPda,  // The vault's public key
    recipient, // The recipient's public key
).Build()

// Get the latest transaction index from the multisig account.
multisig, err := s.MultisigAccount(context.Background())
if err != nil {
    // Handle error
}
transactionIndex := multisig.TransactionIndex + 1

// Create the transaction.
tx, err := s.VaultTransactionCreateTx(
    context.Background(),
    signer,
    vaultIndex,
    transactionIndex,
    []solana.Instruction{vaultInstruction},
)
if err != nil {
    // Handle error
}

// Sign and send the transaction...
```

### Create and Approve a Proposal

To execute a transaction, you first need to create a proposal and have it approved by the required number of members.

#### 1. Create a Proposal

```go
// ...

// Get the latest transaction index.
multisig, err := s.MultisigAccount(context.Background())
if err != nil {
    // Handle error
}
transactionIndex := multisig.TransactionIndex

// Create the proposal transaction.
tx, err := s.ProposalCreateTx(context.Background(), signer, transactionIndex)
if err != nil {
    // Handle error
}

// Sign and send the transaction...
```

#### 2. Approve the Proposal

```go
// ...

// The transaction index of the proposal to approve.
transactionIndex := // The index from the previous step

// Create the approval transaction.
tx, err := s.ProposalApproveTx(context.Background(), signer, transactionIndex)
if err != nil {
    // Handle error
}

// Sign and send the transaction...
```

### Execute a Transaction

Once a proposal is approved, you can execute the transaction.

```go
// ...

// The transaction index of the approved proposal.
transactionIndex := // The index of the approved proposal

// Create the execution transaction.
tx, err := s.VaultTransactionExecuteTx(context.Background(), signer, transactionIndex)
if err != nil {
    // Handle error
}

// Sign and send the transaction...
```

## API Reference

For a complete list of available functions and types, please refer to the [Go Reference](https://pkg.go.dev/github.com/Lee0x273/go-squads).

## Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
