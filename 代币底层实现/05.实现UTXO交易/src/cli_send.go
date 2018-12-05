package main 
import (
	"fmt"
	"log"
)
func (cli *CLI) Send(from,to string,amout int){
	if !ValidateAddress([]byte(from)) {
		log.Panic("ERROR: 发送人的地址错误")
	}
	if !ValidateAddress([]byte(to)) {
		log.Panic("ERROR: 接收人的地址错误")
	}

	//创建区块链
	bc := NewBlockchain()
	defer bc.db.Close()
	//创建utxo集合
	utxo := UTXOSet{bc}
	transaction := NewUTXOTransaction(from,to,amout,&utxo)
	cbtx := NewCoinbaseTX(from,"")
	txs := []*Transaction{transaction,cbtx}
	//挖矿
	block := bc.MineBlock(txs)
	utxo.Update(block)
	//utxo.Reindex()
	fmt.Println("Success~!")
}