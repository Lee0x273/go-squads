package squads

import (
	"context"
	"fmt"
	"squads/generated/squads_multisig_program"
	"testing"

	"github.com/axengine/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func Test_Multisig(t *testing.T) {
	client := rpc.New("http://47.241.179.122:8001/")

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	squads := NewSQuard(multisigPda, client)
	multisig, err := squads.Multisig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.JsonPretty(multisig))
}

func Test_VaultTransactionCreate(t *testing.T) {
	client := rpc.New("http://47.241.179.122:8001/")

	signer, _ := solana.PrivateKeyFromBase58("5vhBxBxp3JXxnMYmkM9q1NwTckCRSiQTHk7aKBes22FRoCQ5QVmTqyTsaLqLy4eDHQJnyb5QPoqzcvKfCiqbYSqH")
	fmt.Println("signer:", signer.PublicKey())

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := NewSQuard(multisigPda, client)
	multisig, err := s.Multisig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	vaultPda, err := s.GetVaultPda(0)
	if err != nil {
		t.Fatal(err)
	}
	vaultInstruction := system.NewTransferInstruction(
		10000000,
		vaultPda,
		signer.PublicKey(),
	).Build()

	transactionIndex := multisig.TransactionIndex + 1
	fmt.Println("transactionIndex=", transactionIndex)

	tx, err := s.CreateVaultTransactionCreate(t.Context(), signer.PublicKey(), 0,
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
	client := rpc.New("http://47.241.179.122:8001/")

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := NewSQuard(multisigPda, client)
	transactionIndex := uint64(1)
	pda, _ := s.GetTransactionPda(uint64(transactionIndex))
	act, err := s.VaultTransactionAccount(t.Context(), pda)
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrettyToStdout(act)

}

func Test_CreateProposal(t *testing.T) {
	client := rpc.New("http://47.241.179.122:8001/")

	signer, _ := solana.PrivateKeyFromBase58("5RKVVjG1RDYD2biSDfNSfakGWEozVjUdgtDQAi74CghZyrCVf4gQot6X2ZBSmmJzMhGhkN9t8hGFdsy2337CbA1E")
	fmt.Println("signer:", signer.PublicKey())

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := NewSQuard(multisigPda, client)
	multisig, err := s.Multisig(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	tx, err := s.CreateProposalCreate(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
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

func Test_ProposalVote(t *testing.T) {
	client := rpc.New("http://47.241.179.122:8001/")

	signer, _ := solana.PrivateKeyFromBase58("5RKVVjG1RDYD2biSDfNSfakGWEozVjUdgtDQAi74CghZyrCVf4gQot6X2ZBSmmJzMhGhkN9t8hGFdsy2337CbA1E")
	fmt.Println("signer:", signer.PublicKey())

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := NewSQuard(multisigPda, client)
	multisig, err := s.Multisig(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	proposalPda, err := s.GetProposalPda(multisig.TransactionIndex)
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := s.ProposalAccount(t.Context(), proposalPda)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.JsonPretty(proposal))

	switch proposal.Status.(type) {
	case *squads_multisig_program.ProposalStatusActive:
		fmt.Println("proposal is active")
	case *squads_multisig_program.ProposalStatusDraft:
		fmt.Println("proposal is draft")
	case *squads_multisig_program.ProposalStatusExecuted:
		fmt.Println("proposal is executed")
	case *squads_multisig_program.ProposalStatusRejected:
		fmt.Println("proposal is rejected")
	default:
		t.Fatal("unknown proposal status")
	}

	tx, err := s.CreateProposalVote(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
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
	client := rpc.New("http://47.241.179.122:8001/")

	signer, _ := solana.PrivateKeyFromBase58("5vhBxBxp3JXxnMYmkM9q1NwTckCRSiQTHk7aKBes22FRoCQ5QVmTqyTsaLqLy4eDHQJnyb5QPoqzcvKfCiqbYSqH")
	fmt.Println("signer:", signer.PublicKey())

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	s := NewSQuard(multisigPda, client)
	multisig, err := s.Multisig(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("multisig.TransactionIndex=", multisig.TransactionIndex)
	// proposalPda, err := s.GetProposalPda(multisig.TransactionIndex)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// proposal, err := s.ProposalAccount(t.Context(), proposalPda)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(utils.JsonPretty(proposal))
	// switch proposal.Status.(type) {
	// case *squads_multisig_program.ProposalStatusActive:
	// 	fmt.Println("proposal is active")
	// case *squads_multisig_program.ProposalStatusDraft:
	// 	fmt.Println("proposal is draft")
	// case *squads_multisig_program.ProposalStatusExecuted:
	// 	fmt.Println("proposal is executed")
	// case *squads_multisig_program.ProposalStatusRejected:
	// 	fmt.Println("proposal is rejected")
	// case *squads_multisig_program.ProposalStatusApproved:
	// 	fmt.Println("proposal is approved")
	// case *squads_multisig_program.ProposalStatusExecuting:
	// 	fmt.Println("proposal is executing")
	// case *squads_multisig_program.ProposalStatusCancelled:
	// 	fmt.Println("proposal is cancelled")
	// default:
	// 	t.Fatal("unknown proposal status")
	// }

	tx, err := s.CreateVaultTransactionExecute(t.Context(), signer.PublicKey(), multisig.TransactionIndex)
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
