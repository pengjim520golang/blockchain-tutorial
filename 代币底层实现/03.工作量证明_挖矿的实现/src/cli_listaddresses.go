package main

import (
	"fmt"
	"log"
)

func (cli *CLI) listAddresses() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	//获取所有的钱包地址
	addresses := wallets.GetAddresses()
	//遍历输出
	for _, address := range addresses {
		fmt.Println(address)
	}
}
