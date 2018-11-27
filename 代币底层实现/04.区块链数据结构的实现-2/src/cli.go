package main 
import (
	"flag"
	"os"
	"fmt"
	"log"
)
//定义命令行客户端
type CLI struct{}
//命令提示帮助方式
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createwallet - 创建一个新的钱包地址")
	fmt.Println("  listaddresses - 遍历输出所有的钱包地址")
	fmt.Println("  printblockchain - 遍历区块链")
	fmt.Println("  createblockchain -address [address] - 创建区块链的创世区块")
	fmt.Println("  mineblock - 挖矿测试")
}
//当用户直接输入 ./bitCoin但没有参数时提示帮助信息
func (cli *CLI) validateArgs(){
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
//运行客户端命令行
func (cli *CLI) Run(){
	cli.validateArgs()
	//注册createwallet命令
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	//注册listaddresses命令
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	//注册printblockchain命令
	printblockchainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)
	//注册createblockchain命令
	createblockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	//给createblockchain命令注册-address选项
	createBlockchainAddress := createblockchainCmd.String("address", "", "address是创世块的奖励地址")
	//注册mineblock命令
	mineblockCmd := flag.NewFlagSet("mineblock", flag.ExitOnError)
	switch os.Args[1] {
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printblockchain":
		err := printblockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createblockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "mineblock":
		err := mineblockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	}
	//如果用户输入的是createwallet就执行对应方法
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	//如果用户输入的是listaddresses就执行对应方法
	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}
	//如果用户输入的是printblockchain就执行对应方法
	if printblockchainCmd.Parsed() {
		cli.printBlockChain()
	}
	//处理createblockchain -address命令
	if createblockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createblockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}
	
	if mineblockCmd.Parsed() {
		cli.Mine()
	}
}