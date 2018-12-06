package main

import (
	"log"
	"fmt"
)

func (cli *CLI) createBlockchain(address string){
	//判断当前钱包地址是否有效
	if !ValidateAddress([]byte(address)) {
		log.Panic("错误的地址:",address)
	}
	bc := CreateBlockChain(address)
	defer bc.db.Close()
	
	//创建utxo_set
	u := UTXOSet{bc}
	u.Reindex()

	fmt.Println("Done !")
}