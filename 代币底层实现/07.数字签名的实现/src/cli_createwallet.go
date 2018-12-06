package main

import "fmt"

func (cli *CLI) createWallet() {
	//创建集合
	wallets, _ := NewWallets()
	//创建钱包并把钱包写入map[string]*Wallet字典中
	address := wallets.CreateWallet()
	//保持集合序列化数据到文件中
	wallets.SaveToFile()

	fmt.Printf("你创建的新钱包地址是: %s\n", address)
}