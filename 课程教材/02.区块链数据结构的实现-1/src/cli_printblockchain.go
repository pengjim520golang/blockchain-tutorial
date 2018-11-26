package main

import (
	"fmt"

)

func (cli *CLI) printBlockChain() {
	//创建创世区块
	txIn := TXInput{[]byte{},-1,nil,[]byte("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks")}
	txOut := TXOutput{50,nil}
	coinbase := &Transaction{[]byte("11111111111"),[]TXInput{txIn},[]TXOutput{txOut}}
	block := NewGenesisBlock(coinbase)
	//创建拥有创世区块的区块链
	bc := NewBlockchain(block)

	//创建交易1
	txIn = TXInput{[]byte("11111111111"),
				   0,
				   []byte("1EEdwyuhhNXyYCheqibbTDnZygZD4ypiop"),
				   []byte("pubkey")}
	txOut = TXOutput{20,
		            []byte("1NBP1H3VMzPKNx6Lknk2pL7dFJTUQswC5n")}
	txOne := &Transaction{[]byte("22222222"), 
						  []TXInput{txIn},
						  []TXOutput{txOut}}
	//将交易信息加入区块链
	bc.AddBlock([]*Transaction{txOne})
	
	//创建交易2
	txIn = TXInput{[]byte("11111111111"),
					0,
					[]byte("1EEdwyuhhNXyYCheqibbTDnZygZD4ypiop"),
					[]byte("pubkey")}
	txOut = TXOutput{10,
		             []byte("1NBP1H3VMzPKNx6Lknk2pL7dFJTUQswC5n")}
	txSec := &Transaction{[]byte("22222222"),
						  []TXInput{txIn},
						  []TXOutput{txOut}}
	//将交易信息加入区块链
	bc.AddBlock([]*Transaction{txSec})

	//遍历区块链
	for _,block := range bc.blocks{
		fmt.Printf("Hash:%x\nPrevBlockHash:%x\n",block.Hash,block.PrevBlockHash)
		fmt.Println()
	} 

}
