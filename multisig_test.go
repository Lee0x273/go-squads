package squads

import (
	"context"
	"fmt"
	"testing"

	"github.com/axengine/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func Test_Multisig(t *testing.T) {
	client := rpc.New(rpc.DevNet.RPC)
	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	squads := New(client, multisigPda)
	multisig, err := squads.MultisigAccount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.JsonPretty(multisig))
}

func Test_VaultTransactionCreate(t *testing.T) {
	client := rpc.New(rpc.DevNet.RPC)

	signer, _ := solana.PrivateKeyFromSolanaKeygenFile("creator.json")
	fmt.Println("signer:", signer.PublicKey())

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := New(client, multisigPda)
	multisig, err := s.MultisigAccount(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	vaultPda, err := GetVaultPda(s.multisigPda, 0)
	if err != nil {
		t.Fatal(err)
	}
	vaultInstruction := system.NewTransferInstruction(
		100000000,
		vaultPda,
		signer.PublicKey(),
	).Build()

	transactionIndex := multisig.TransactionIndex + 1
	fmt.Println("transactionIndex=", transactionIndex)

	tx, err := s.CreateVaultTransactionCreateTx(t.Context(), signer.PublicKey(), 0,
		transactionIndex,
		[]solana.Instruction{vaultInstruction})
	if err != nil {
		t.Fatal(err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if signer.PublicKey().Equals(key) {
			return &signer
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("transaction broadcasted,signature: %s\n", sig.String())
}

func Test_GetVaultTransaction(t *testing.T) {
	client := rpc.New(rpc.DevNet.RPC)

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := New(client, multisigPda)
	transactionIndex := uint64(22)
	pda, _ := GetTransactionPda(s.multisigPda, uint64(transactionIndex))
	act, err := s.VaultTransactionAccount(t.Context(), pda)
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrettyToStdout(act)

}

func Test_CreateProposal(t *testing.T) {
	client := rpc.New(rpc.DevNet.RPC)

	signer, _ := solana.PrivateKeyFromSolanaKeygenFile("creator.json")
	fmt.Println("signer:", signer.PublicKey())

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := New(client, multisigPda)
	multisig, err := s.MultisigAccount(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	tx, err := s.CreateProposalCreateTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if signer.PublicKey().Equals(key) {
			return &signer
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("transaction broadcasted,signature: %s\n", sig.String())
}

func Test_GetProposal(t *testing.T) {
	client := rpc.New(rpc.DevNet.RPC)

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := New(client, multisigPda)
	transactionIndex := uint64(22)
	pda, _ := GetProposalPda(s.multisigPda, transactionIndex)
	act, err := s.ProposalAccount(t.Context(), pda)
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrettyToStdout(act)

}

func Test_ProposalApprove(t *testing.T) {
	client := rpc.New(rpc.DevNet.RPC)

	signer, _ := solana.PrivateKeyFromSolanaKeygenFile("creator.json") // or secondmember.json
	fmt.Println("signer:", signer.PublicKey())

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := New(client, multisigPda)
	multisig, err := s.MultisigAccount(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	proposalPda, err := GetProposalPda(s.multisigPda, multisig.TransactionIndex)
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := s.ProposalAccount(t.Context(), proposalPda)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.JsonPretty(proposal))

	fmt.Println(GetProposalStatus(proposal.Status))

	tx, err := s.CreateProposalApproveTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if signer.PublicKey().Equals(key) {
			return &signer
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("transaction broadcasted,signature: %s\n", sig.String())
}

func Test_VaultTransactionExecute(t *testing.T) {
	client := rpc.New(rpc.DevNet.RPC)

	signer, _ := solana.PrivateKeyFromSolanaKeygenFile("creator.json")
	fmt.Println("signer:", signer.PublicKey())

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := New(client, multisigPda)
	multisig, err := s.MultisigAccount(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("multisig.TransactionIndex=", multisig.TransactionIndex)

	tx, err := s.CreateVaultTransactionExecuteTx(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if signer.PublicKey().Equals(key) {
			return &signer
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("transaction broadcasted,signature: %s\n", sig.String())
}
