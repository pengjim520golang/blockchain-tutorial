package main 

import (
	"fmt"
	"crypto/elliptic"
	"bytes"
	"encoding/gob"
	"log"
	"io/ioutil"
	"os"
)

const WalletsFile = "wallets.dat"
type Wallets struct {
	Wallets map[string]*Wallet 
}

func NewWallets() (*Wallets,error){
	wallets := &Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	return wallets,err 
}
//创建一个钱包,并把钱包加入到集合中
func (wallets *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s",wallet.GetAddress())
	wallets.Wallets[address] = wallet
	return address
}

//从文件中获取钱包集合
func (ws *Wallets) LoadFromFile() error {
	if _,err:=os.Stat(WalletsFile);os.IsNotExist(err){
		return err;
	}
	//读取文件中的序列化字节数据
	contents , err := ioutil.ReadFile(WalletsFile)
	if err != nil {
		log.Panic(err)
	}
	//反序列
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(contents))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	//构造钱包集合
	ws.Wallets = wallets.Wallets
	return nil 
}
//保存钱包集合的序列化数据
func (ws Wallets) SaveToFile(){
	var buffer bytes.Buffer
	//这代码只是一个标示，没什么意思，只是一个良好的编程习惯,不写也行
	gob.Register(elliptic.P256())
	//序列化
	encoder := gob.NewEncoder( &buffer )
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	//把序列化数据写入文件
	err = ioutil.WriteFile(WalletsFile,buffer.Bytes(),0666)
	if err != nil {
		log.Panic(err)
	}
}
//根据地址获取一个钱包实例
func (wallets Wallets) GetWallet(address string) Wallet{
	//这要非常注意，因为Wallets中的map返回的是*Wallet
	return *wallets.Wallets[address]
}
//获取所有都钱包地址
func (wallets *Wallets) GetAddresses() []string{
	var addresses []string 

	for address := range wallets.Wallets{
		addresses = append(addresses,address)
	}

	return addresses
}

