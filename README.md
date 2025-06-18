# golang sdk of squads v4
- support `SmallVec`

## usage
- new squads
```
	import(
		"github.com/Lee0x273/go-squads"
	)
	client := rpc.New(rpc.DevNet.RPC)
	multisigPda := solana.MustPublicKeyFromBase58("...")
	s := squads.New(client, multisigPda)
```
- create vault transaction
```
	tx, err := s.CreateVaultTransactionCreateTx(t.Context(), signer.PublicKey(), vaultIndex,
		transactionIndex,
		[]solana.Instruction{vaultInstruction})
```
- create proposal transaction
```
	tx, err := s.CreateProposalCreateTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
```
- create approve transaction
```
	tx, err := s.CreateProposalApproveTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
```
- create execute transaction
```
	tx, err := s.CreateVaultTransactionExecuteTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
```

see more: squads_test.go