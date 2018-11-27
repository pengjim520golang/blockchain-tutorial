package main

import (
	"strconv"
	"fmt"
)

func (cli *CLI) printBlockChain() {
	bc := NewBlockchain()
	it := bc.Iterator()
	for {
		block := it.Next()
		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))		

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
