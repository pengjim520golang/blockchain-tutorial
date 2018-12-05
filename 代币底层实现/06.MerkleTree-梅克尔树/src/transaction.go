package main

import (
	"crypto/rand"
	"encoding/hex"
	"crypto/sha256"
	"log"
	"encoding/gob"
	"bytes"
	"fmt"
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
		//data = "coinbase transaction"
		randData := make([]byte,20)
		_,err :=  rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}
	//创建一笔输入
	txIn := TXInput{[]byte{},-1,nil,[]byte(data)}
	txOut := NewTXOutput(subsidy,to) 
	cbtx := Transaction{nil,[]TXInput{txIn},[]TXOutput{*txOut}}
	cbtx.ID = cbtx.Hash()
	return &cbtx 
}
//判断当前交易是否属于Coinbase
func (tx *Transaction) IsCoinbase() bool{
	return (len(tx.Vint)==1 && tx.Vint[0].Vout==-1 && len(tx.Vint[0].Txid)==0)
}
//创建一笔交易
func NewUTXOTransaction(from , to string,amount int,utxoSet *UTXOSet) *Transaction{
	var Txinputs []TXInput
	var Txoutputs []TXOutput
	//获取钱包中的PubKey
	wallets,err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc,outputs := utxoSet.FindSpendableOutputs(pubKeyHash,amount)
	if acc < amount {
		log.Panic("没有足够完成本次交易的金额")
	}
	for txid,outputIndexes := range outputs{
		txID,err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		//遍历可用的输出索引
		for _,outIndex := range outputIndexes{
			txIn := TXInput{txID,outIndex,nil,wallet.PublicKey}
			Txinputs = append(Txinputs,txIn)
		}
	}
	//创建一笔交易输出
	txOut := *NewTXOutput(amount,to)
	Txoutputs = append(Txoutputs,txOut)
	if acc > amount {
		change := acc - amount 
		//创建一逼找零输出
		changeOut := *NewTXOutput(change,from)
		Txoutputs = append(Txoutputs,changeOut)
	}
	transaction := Transaction{nil,Txinputs,Txoutputs}
	transaction.ID = transaction.Hash()
	return &transaction
}