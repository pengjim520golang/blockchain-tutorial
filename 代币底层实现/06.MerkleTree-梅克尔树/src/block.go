package main

import (
	"log"
	"encoding/gob"
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
	block := &Block{time.Now().Unix(),transactions,prevBlockHash,[]byte{},0}

	//加入工作量证明
	pow := NewProofOfWork(block)
	//进行挖矿
	nonce,hash := pow.Run()
	//设置当前区块的Nonce随机数
	block.Nonce = nonce
	//设置区块Hash
	block.Hash = hash[:]
	return block
}
//创建创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
//把交易打包方便在区块中存储
func (block *Block) HashTransactions() []byte{
	var transactions [][]byte 
	for _,tx := range block.Transactions{
		transactions = append(transactions,tx.Serialize())
	}
	merkleTree := NewMerkleTree(transactions)

	return merkleTree.RootNode.Data
}

//序列化区块
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

//反序列化区块
func UnSerializeBlock(blockBytes []byte) *Block{
	var block Block
	decoder := gob.NewDecoder( bytes.NewReader(blockBytes) )
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic("err is ",err)
	}
	return &block	
}