package main

func (cli *CLI) Mine() {
	//创建区块链实例
	bc := NewBlockchain()
	//创建交易1
	txIn := TXInput{[]byte("11111111111"),0,[]byte("1EEdwyuhhNXyYCheqibbTDnZygZD4ypiop"),[]byte("pubkey")}
	txOut := TXOutput{20,[]byte("1NBP1H3VMzPKNx6Lknk2pL7dFJTUQswC5n")}
	txOne := &Transaction{[]byte("22222222"), []TXInput{txIn},[]TXOutput{txOut}}
	bc.MineBlock([]*Transaction{txOne})
	//创建交易2
	txIn = TXInput{[]byte("11111111111"),0,[]byte("1EEdwyuhhNXyYCheqibbTDnZygZD4ypiop"),[]byte("pubkey")}
	txOut = TXOutput{10,[]byte("1NBP1H3VMzPKNx6Lknk2pL7dFJTUQswC5n")}
	txSec := &Transaction{[]byte("22222222"),[]TXInput{txIn},[]TXOutput{txOut}}
	bc.MineBlock([]*Transaction{txSec})
}