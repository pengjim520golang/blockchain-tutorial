package main

import (
	"crypto/sha256"
	"bytes"
	"time"
)

//区块数据结构
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

//创建一个区块实例
func NewBlock(transactions []*Transaction,prevBlockHash []byte) *Block{
	block := &Block{time.Now().Unix(),
					transactions,
					prevBlockHash,
					[]byte{},
					0}
	//创建区块地址
	block.setHash()
	return block
}
//创建创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
//把交易打包方便在区块中存储(在比特币中应该把交易打包成为一个梅克尔树,本小节是简单实现)
func (block *Block) HashTransactions() []byte{
	var transactions [][]byte 
	for _,tx := range block.Transactions{
		transactions = append(transactions,tx.Serialize())
	}
	txHash := sha256.Sum256( bytes.Join(transactions,[]byte{}) )
	return txHash[:]
}
//设置Block的Hash(这个工作我们应该在工作量(挖矿)中实现,本小节是简单实现)
func (block *Block) setHash(){
	//组合用于哈希的区块数据
	hashData := bytes.Join([][]byte{ IntToHex(block.Timestamp),
									 IntToHex(int64(block.Nonce)),
									 block.HashTransactions(),
									 block.PrevBlockHash},[]byte{})

	var hash [32]byte 
	hash = sha256.Sum256(hashData)
	//设置当前区块的哈希地址
	block.Hash = hash[:]
}