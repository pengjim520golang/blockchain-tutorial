package main

import (
	"crypto/sha256"
	"log"
	"encoding/gob"
	"bytes"
)
//基币奖励
const subsidy = 10

type Transaction struct{
	ID []byte
	Vint []TXInput
	Vout []TXOutput
}
//序列化交易
func (tx Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
//把当前交易进行哈希
func (tx *Transaction) Hash() []byte{
	var hash [32]byte
	txCopy := *tx 
	txCopy.ID = []byte{}
	hash = sha256.Sum256(tx.Serialize()) 
	return hash[:]
}
//创建coinbase交易
func NewCoinbaseTX( to,data string ) *Transaction{
	if data =="" {
		data = "transaction"
	}
	//定义coinbase交易
	txIn := TXInput{[]byte{},-1,nil,[]byte(data)}
	txOut := TXOutput{subsidy,[]byte(to)}
	cbtx := Transaction{[]byte{},[]TXInput{txIn},[]TXOutput{txOut}}
	cbtx.ID = cbtx.Hash()
	return &cbtx 
}

