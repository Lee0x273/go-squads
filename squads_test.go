package squads

import (
	"context"
	"fmt"
	"squads/generated/squads_multisig_program"
	"testing"

	"github.com/axengine/utils"
	"github.com/davecgh/go-spew/spew"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func TestXxx(t *testing.T) {
	client := rpc.New("http://47.241.179.122:8001/")

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")
	_ = multisigPda

	out, err := client.GetAccountInfo(context.Background(), multisigPda)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.JsonPretty(out))

	data := out.Value.Data.GetBinary()

	// 打印原始数据用于调试
	fmt.Printf("Raw data length: %d\n", len(data))
	fmt.Printf("Discriminator: %x\n", data[:8])
	fmt.Printf("Raw data: %x\n", data)

	multisig := &squads_multisig_program.Multisig{}
	decoder := ag_binary.NewBorshDecoder(data)
	if err := multisig.UnmarshalWithDecoder(decoder); err != nil {
		t.Fatal(err)
	}

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
	//return
	if err := vaultTransactionCreate(client); err != nil {
		t.Fatal(err)
	}
}

func vaultTransactionCreate(client *rpc.Client) error {
	signer, _ := solana.PrivateKeyFromBase58("5vhBxBxp3JXxnMYmkM9q1NwTckCRSiQTHk7aKBes22FRoCQ5QVmTqyTsaLqLy4eDHQJnyb5QPoqzcvKfCiqbYSqH")

	multisigPda := solana.MustPublicKeyFromBase58("G26QSXWEdY11iue8Dw2aushtw7hhVF5zHDhSXqSJGRLA")

	vaultIndex := uint8(0)
	amount := uint64(10000000) // 0.01
	// 获取 vault PDA
	vaultPda, err := NewSQuard().GetVaultPda(multisigPda, vaultIndex)
	if err != nil {
		return fmt.Errorf("获取 vault pda 失败: %v", err)
	}
	fmt.Println("vault pda:", vaultPda.String())

	// 获取最新区块哈希用于构造TransactionMessage
	recent, err := client.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("获取最新区块哈希失败: %v", err)
	}

	// 获取当前multisig的交易索引
	out, err := client.GetAccountInfo(context.Background(), multisigPda)
	if err != nil {
		return fmt.Errorf("获取multisig账户信息失败: %v", err)
	}
	data := out.Value.Data.GetBinary()
	multisigInfo := &squads_multisig_program.Multisig{}
	decoder := ag_binary.NewBorshDecoder(data)
	if err := multisigInfo.UnmarshalWithDecoder(decoder); err != nil {
		return fmt.Errorf("解析multisig账户信息失败: %v", err)
	}

	// 使用当前交易索引+1作为新交易索引
	transactionIndex := multisigInfo.TransactionIndex + 1
	fmt.Printf("当前交易索引: %d, 新交易索引: %d\n", multisigInfo.TransactionIndex, transactionIndex)

	// 获取交易 PDA
	transactionPda, err := NewSQuard().GetTransactionPda(multisigPda, transactionIndex)
	if err != nil {
		return fmt.Errorf("获取 transaction pda 失败: %v", err)
	}
	fmt.Println("transaction pda:", transactionPda.String())

	// transfer lamports from vaultPda to other
	vaultTx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				vaultPda,
				signer.PublicKey(),
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(vaultPda),
	)
	if err != nil {
		panic(err)
	}

	txBuf, _ := vaultTx.Message.MarshalV0() //??

	args := squads_multisig_program.VaultTransactionCreateArgs{
		VaultIndex:         vaultIndex,
		EphemeralSigners:   0,
		TransactionMessage: txBuf, //??
	}

	instruction := squads_multisig_program.NewVaultTransactionCreateInstruction(
		args,
		multisigPda,
		transactionPda,
		signer.PublicKey(), // creator
		signer.PublicKey(), // feePayer
		system.ProgramID,
	)

	tx, _ := solana.NewTransaction(
		[]solana.Instruction{instruction.Build()},
		recent.Value.Blockhash,
		solana.TransactionPayer(signer.PublicKey()),
	)
	// tx.Sign(...)

	// 签名交易
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if signer.PublicKey().Equals(key) {
			return &signer
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("签名交易失败: %v", err)
	}

	// 发送交易
	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		return fmt.Errorf("发送交易失败: %v", err)
	}

	fmt.Printf("交易已发送，签名: %s\n", sig.String())
	return nil
}

func TestTokenSupply(t *testing.T) {
	endpoint := rpc.MainNetBeta_RPC
	client := rpc.New(endpoint)

	pubKey := solana.MustPublicKeyFromBase58("YBTCKqqMriwfFXVsuem82bkGFjigmUWFWD28JgtrumX") // serum token
	out, err := client.GetTokenSupply(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}
	spew.Dump(out)
}
