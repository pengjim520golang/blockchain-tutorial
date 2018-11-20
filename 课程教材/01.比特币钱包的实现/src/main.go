package main

import (
	"fmt"
)

func main(){
	/*
	wallet := NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("钱包地址:%s\n",address)
	fmt.Println("有效性:",ValidateAddress(address)) */

	wallets , _ := NewWallets()
	walletAdr := wallets.CreateWallet()
	fmt.Println("新建钱包：",walletAdr)
	wallets.SaveToFile()
	addresses := wallets.GetAddresses()
	fmt.Println("钱包集合如下：")
	for _,address := range addresses{
		fmt.Println(address)
	}
}