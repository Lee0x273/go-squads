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
	client := rpc.New("http://47.241.179.122:8001/")

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	squads := NewSQuard(multisigPda, client)
	multisig, err := squads.Multisig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.JsonPretty(multisig))

	// 在输出信息时添加 discriminator
	// fmt.Printf("Discriminator: %x\n", multisig.Discriminator)
	fmt.Printf("多签钱包地址: %s\n", multisigPda.String())
	fmt.Printf("创建密钥: %s\n", multisig.CreateKey.String())
	fmt.Printf("配置权限: %s\n", multisig.ConfigAuthority.String())
	fmt.Printf("签名阈值: %d\n", multisig.Threshold)
	fmt.Printf("时间锁: %d秒\n", multisig.TimeLock)
	fmt.Printf("交易索引: %d\n", multisig.TransactionIndex)
	fmt.Printf("过期交易索引: %d\n", multisig.StaleTransactionIndex)
	if multisig.RentCollector != nil {
		fmt.Printf("租金收集者: %s\n", multisig.RentCollector.String())
	}
	fmt.Printf("Bump: %d\n", multisig.Bump)
	// fmt.Printf("成员数量: %d\n", multisig.MemberCount)
	fmt.Println("成员列表:")
	for i, member := range multisig.Members {
		fmt.Printf("  %d. %s %d\n", i+1, member.Key.String(), member.Permissions.Mask)
	}
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

	tx, err := s.CreateVaultTransactionCreate(t.Context(), signer.PublicKey(), 0,
		multisig.TransactionIndex+1,
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
