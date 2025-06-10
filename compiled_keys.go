package squads

import (
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
)

// CompiledKeyMeta 对应TypeScript中的CompiledKeyMeta
type CompiledKeyMeta struct {
	IsSigner   bool `json:"isSigner"`
	IsWritable bool `json:"isWritable"`
	IsInvoked  bool `json:"isInvoked"`
}

// CompiledKeys 对应TypeScript中的CompiledKeys类
type CompiledKeys struct {
	Payer      solana.PublicKey           `json:"payer"`
	KeyMetaMap map[string]CompiledKeyMeta `json:"keyMetaMap"`
}

// AccountKeysFromLookups 从查找表获取的账户密钥
type AccountKeysFromLookups struct {
	Writable []solana.PublicKey `json:"writable"`
	Readonly []solana.PublicKey `json:"readonly"`
}

// MessageV0 对应Solana V0消息
type MessageV0 struct {
	Header               solana.MessageHeader               `json:"header"`
	StaticAccountKeys    []solana.PublicKey                 `json:"staticAccountKeys"`
	RecentBlockhash      solana.Hash                        `json:"recentBlockhash"`
	CompiledInstructions []solana.CompiledInstruction       `json:"compiledInstructions"`
	AddressTableLookups  []solana.MessageAddressTableLookup `json:"addressTableLookups"`
}

// MessageAccountKeys 消息账户密钥管理
type MessageAccountKeys struct {
	StaticAccountKeys      []solana.PublicKey     `json:"staticAccountKeys"`
	AccountKeysFromLookups AccountKeysFromLookups `json:"accountKeysFromLookups"`
}

// NewCompiledKeys 创建新的CompiledKeys实例
func NewCompiledKeys(payer solana.PublicKey, keyMetaMap map[string]CompiledKeyMeta) *CompiledKeys {
	return &CompiledKeys{
		Payer:      payer,
		KeyMetaMap: keyMetaMap,
	}
}

// Compile 编译指令和付款人密钥
func CompileKeys(instructions []solana.Instruction, payer solana.PublicKey) *CompiledKeys {
	keyMetaMap := make(map[string]CompiledKeyMeta)

	getOrInsertDefault := func(pubkey solana.PublicKey) *CompiledKeyMeta {
		address := pubkey.String()
		if keyMeta, exists := keyMetaMap[address]; exists {
			return &keyMeta
		}

		keyMeta := CompiledKeyMeta{
			IsSigner:   false,
			IsWritable: false,
			IsInvoked:  false,
		}
		keyMetaMap[address] = keyMeta
		return &keyMeta
	}

	// 设置付款人为签名者和可写
	payerKeyMeta := getOrInsertDefault(payer)
	payerKeyMeta.IsSigner = true
	payerKeyMeta.IsWritable = true
	keyMetaMap[payer.String()] = *payerKeyMeta

	// 处理指令
	for _, ix := range instructions {
		// 在multisig版本中，程序ID不标记为调用
		programKeyMeta := getOrInsertDefault(ix.ProgramID())
		programKeyMeta.IsInvoked = false
		keyMetaMap[ix.ProgramID().String()] = *programKeyMeta

		// 处理账户密钥
		for _, accountMeta := range ix.Accounts() {
			keyMeta := getOrInsertDefault(accountMeta.PublicKey)
			keyMeta.IsSigner = keyMeta.IsSigner || accountMeta.IsSigner
			keyMeta.IsWritable = keyMeta.IsWritable || accountMeta.IsWritable
			keyMetaMap[accountMeta.PublicKey.String()] = *keyMeta
		}
	}

	return NewCompiledKeys(payer, keyMetaMap)
}

// GetMessageComponents 获取消息组件
func (ck *CompiledKeys) GetMessageComponents() (solana.MessageHeader, []solana.PublicKey) {
	var writableSigners, readonlySigners, writableNonSigners, readonlyNonSigners []string

	for address, meta := range ck.KeyMetaMap {
		if meta.IsSigner && meta.IsWritable {
			writableSigners = append(writableSigners, address)
		} else if meta.IsSigner && !meta.IsWritable {
			readonlySigners = append(readonlySigners, address)
		} else if !meta.IsSigner && meta.IsWritable {
			writableNonSigners = append(writableNonSigners, address)
		} else {
			readonlyNonSigners = append(readonlyNonSigners, address)
		}
	}

	header := solana.MessageHeader{
		NumRequiredSignatures:       uint8(len(writableSigners) + len(readonlySigners)),
		NumReadonlySignedAccounts:   uint8(len(readonlySigners)),
		NumReadonlyUnsignedAccounts: uint8(len(readonlyNonSigners)),
	}

	// 构建静态账户密钥数组
	var staticAccountKeys []solana.PublicKey

	// 添加可写签名者
	for _, address := range writableSigners {
		pubkey, _ := solana.PublicKeyFromBase58(address)
		staticAccountKeys = append(staticAccountKeys, pubkey)
	}

	// 添加只读签名者
	for _, address := range readonlySigners {
		pubkey, _ := solana.PublicKeyFromBase58(address)
		staticAccountKeys = append(staticAccountKeys, pubkey)
	}

	// 添加可写非签名者
	for _, address := range writableNonSigners {
		pubkey, _ := solana.PublicKeyFromBase58(address)
		staticAccountKeys = append(staticAccountKeys, pubkey)
	}

	// 添加只读非签名者
	for _, address := range readonlyNonSigners {
		pubkey, _ := solana.PublicKeyFromBase58(address)
		staticAccountKeys = append(staticAccountKeys, pubkey)
	}

	return header, staticAccountKeys
}

// ExtractTableLookup 提取表查找
func (ck *CompiledKeys) ExtractTableLookup(lookupTable addresslookuptable.KeyedAddressLookupTable) (*solana.MessageAddressTableLookup, *AccountKeysFromLookups, bool) {
	writableIndexes, drainedWritableKeys := ck.drainKeysFoundInLookupTable(
		lookupTable.State.Addresses,
		func(keyMeta CompiledKeyMeta) bool {
			return !keyMeta.IsSigner && !keyMeta.IsInvoked && keyMeta.IsWritable
		},
	)

	readonlyIndexes, drainedReadonlyKeys := ck.drainKeysFoundInLookupTable(
		lookupTable.State.Addresses,
		func(keyMeta CompiledKeyMeta) bool {
			return !keyMeta.IsSigner && !keyMeta.IsInvoked && !keyMeta.IsWritable
		},
	)

	// 如果没有找到密钥，不提取查找
	if len(writableIndexes) == 0 && len(readonlyIndexes) == 0 {
		return nil, nil, false
	}

	return &solana.MessageAddressTableLookup{
			AccountKey:      lookupTable.Key,
			WritableIndexes: writableIndexes,
			ReadonlyIndexes: readonlyIndexes,
		},
		&AccountKeysFromLookups{
			Writable: drainedWritableKeys,
			Readonly: drainedReadonlyKeys,
		},
		true
}

// drainKeysFoundInLookupTable 从查找表中提取找到的密钥
func (ck *CompiledKeys) drainKeysFoundInLookupTable(lookupTableEntries []solana.PublicKey, keyMetaFilter func(CompiledKeyMeta) bool) ([]uint8, []solana.PublicKey) {
	var lookupTableIndexes []uint8
	var drainedKeys []solana.PublicKey

	for address, keyMeta := range ck.KeyMetaMap {
		if keyMetaFilter(keyMeta) {
			key, _ := solana.PublicKeyFromBase58(address)

			// 查找密钥在查找表中的索引
			for i, entry := range lookupTableEntries {
				if entry.Equals(key) {
					lookupTableIndexes = append(lookupTableIndexes, uint8(i))
					drainedKeys = append(drainedKeys, key)
					delete(ck.KeyMetaMap, address)
					break
				}
			}
		}
	}

	return lookupTableIndexes, drainedKeys
}

// CompileInstructions 编译指令
func (mk *MessageAccountKeys) CompileInstructions(instructions []solana.Instruction) []solana.CompiledInstruction {
	// 创建账户索引映射
	accountIndexMap := make(map[string]uint16)
	index := uint16(0)

	// 添加静态账户密钥
	for _, key := range mk.StaticAccountKeys {
		accountIndexMap[key.String()] = index
		index++
	}

	// 添加查找表中的可写账户
	for _, key := range mk.AccountKeysFromLookups.Writable {
		accountIndexMap[key.String()] = index
		index++
	}

	// 添加查找表中的只读账户
	for _, key := range mk.AccountKeysFromLookups.Readonly {
		accountIndexMap[key.String()] = index
		index++
	}

	var compiledInstructions []solana.CompiledInstruction

	for _, instruction := range instructions {
		// 获取程序ID索引
		programIDIndex := accountIndexMap[instruction.ProgramID().String()]

		// 编译账户索引
		var accounts []uint16
		for _, accountMeta := range instruction.Accounts() {
			accountIndex := accountIndexMap[accountMeta.PublicKey.String()]
			accounts = append(accounts, accountIndex)
		}

		instructionData, _ := instruction.Data()
		compiledInstructions = append(compiledInstructions, solana.CompiledInstruction{
			ProgramIDIndex: programIDIndex,
			Accounts:       accounts,
			Data:           instructionData,
		})
	}

	return compiledInstructions
}

// CompileToWrappedMessageV0 编译到包装的V0消息
func CompileToWrappedMessageV0(payerKey solana.PublicKey,
	recentBlockhash solana.Hash,
	instructions []solana.Instruction,
	addressLookupTableAccounts []addresslookuptable.KeyedAddressLookupTable) *MessageV0 {

	// 编译密钥
	compiledKeys := CompileKeys(instructions, payerKey)

	// 初始化地址表查找和查找表账户密钥
	var addressTableLookups []solana.MessageAddressTableLookup
	accountKeysFromLookups := AccountKeysFromLookups{
		Writable: []solana.PublicKey{},
		Readonly: []solana.PublicKey{},
	}

	// 处理地址查找表账户
	for _, lookupTable := range addressLookupTableAccounts {
		if lookup, keys, found := compiledKeys.ExtractTableLookup(lookupTable); found {
			addressTableLookups = append(addressTableLookups, *lookup)
			accountKeysFromLookups.Writable = append(accountKeysFromLookups.Writable, keys.Writable...)
			accountKeysFromLookups.Readonly = append(accountKeysFromLookups.Readonly, keys.Readonly...)
		}
	}

	// 获取消息组件
	header, staticAccountKeys := compiledKeys.GetMessageComponents()

	// 创建消息账户密钥管理器
	accountKeys := &MessageAccountKeys{
		StaticAccountKeys:      staticAccountKeys,
		AccountKeysFromLookups: accountKeysFromLookups,
	}

	// 编译指令
	compiledInstructions := accountKeys.CompileInstructions(instructions)

	// 创建并返回V0消息
	return &MessageV0{
		Header:               header,
		StaticAccountKeys:    staticAccountKeys,
		RecentBlockhash:      recentBlockhash,
		CompiledInstructions: compiledInstructions,
		AddressTableLookups:  addressTableLookups,
	}
}
