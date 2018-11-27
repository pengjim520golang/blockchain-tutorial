package main

import (
	"fmt"
)

func (cli *CLI) createBlockchain(address string){
	bc := CreateBlockChain(address)
	defer bc.db.Close()
	fmt.Println("Done !")
}