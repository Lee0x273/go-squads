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
	tx, err := s.VaultTransactionCreateTx(t.Context(), signer.PublicKey(), vaultIndex,
		transactionIndex,
		[]solana.Instruction{vaultInstruction})
```
- create proposal transaction
```
	tx, err := s.ProposalCreateTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
```
- create approve transaction
```
	tx, err := s.ProposalApproveTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
```
- create execute transaction
```
	tx, err := s.VaultTransactionExecuteTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
```

see more: multisig_test.go