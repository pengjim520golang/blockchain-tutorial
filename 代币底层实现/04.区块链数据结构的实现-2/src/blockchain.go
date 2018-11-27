package main 

import (
	"fmt"
	"log"
	"os"
	"github.com/boltdb/bolt"
)

const blockDBFile = "blockchain.db"
const blockTable = "blocks"
const baseInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
/** 
实现持久化的区块链使用bolt DB本地文件kv数据库
tip 属性用于指向区块链最新的元素
db 表示保存区块的数据库文件
**/
type Blockchain struct {
	tip []byte 
	db *bolt.DB
}
//判断区块数据库是否已经创建,如果已经创建则代表已经拥有了创世区块
func dbExists() bool {
	if _,err := os.Stat(blockDBFile); os.IsNotExist(err) {
		return false  
	}
	return true 
}
//使用某个地址作为创世区块的矿工地址,创世区块只能被创建一次
func CreateBlockChain(address string) *Blockchain{
	if dbExists(){
		fmt.Println("创世区块已经建立...")
		os.Exit(1)
	}
	//记录创世区块的hash
	var tip []byte 
	//建立基币交易,这是一笔凭空得出来的奖励
	cbtx := NewCoinbaseTX(address,baseInfo)
	//建立创世区块
	block := NewGenesisBlock(cbtx)
	//创建数据库
	db,err := bolt.Open(blockDBFile,0666,nil)
	if err != nil {
		log.Panic(err)
	}
	//更新区块到bolt db中
	err = db.Update(func(tx *bolt.Tx) error {
		//创建表
		bucket,err := tx.CreateBucket( []byte(blockTable) )
		if err != nil {
			log.Panic(err)
		}
		// hash -> block serialize
		err = bucket.Put(block.Hash,block.Serialize())
		if err != nil {
			log.Panic(err)
		}
		// l 表示最新的区块哈希
		err = bucket.Put([]byte("l"),block.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = block.Hash
		return nil  
	})
	if err != nil {
		log.Panic(err)
	}
	//创建含有创世区块的区块链
	bc := Blockchain{tip,db}
	return &bc 
}
//建立区块链实例
func NewBlockchain() *Blockchain{
	if dbExists() == false {
		fmt.Println("必须创建创世区块")
		os.Exit(1)
	}
	var tip []byte 
	db,err :=  bolt.Open(blockDBFile,0666,nil)
	if err != nil {
		log.Panic(err)	
	}
	err = db.Update(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte(blockTable))
		//获取区块链中最新的区块哈希
		tip = bucket.Get([]byte("l"))
		return nil 
	})
	if err != nil {
		log.Panic(err)
	}
	//创建区块链实例
	bc := Blockchain{tip,db}
	return &bc 
}

//挖矿
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block{
	var lastHash []byte 
	err := bc.db.View( func(tx *bolt.Tx) error{
		bucket := tx.Bucket( []byte(blockTable) )
		lastHash = bucket.Get( []byte("l") )
		return nil 
	})
	if err != nil {
		log.Panic(err)
	}
	//创建一个新对区块
	block := NewBlock(transactions,lastHash) 
	//把区块加入到区块链中
	err = bc.db.Update(func(tx *bolt.Tx) error{
		bucket := tx.Bucket( []byte(blockTable) )
		
		err := bucket.Put(block.Hash,block.Serialize())
		if err != nil {
			log.Panic(err)
		}
		//更新 l 为当前新区块的哈希
		err = bucket.Put([]byte("l"),block.Hash)
		if err != nil {
			log.Panic(err)
		}
		//更新区块链的tip
		bc.tip = block.Hash
		return nil 
	})
	if err != nil {
		log.Panic(err)
	}	
	return block
}

//获取区块链迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator{
	return &BlockchainIterator{bc.tip,bc.db}
}