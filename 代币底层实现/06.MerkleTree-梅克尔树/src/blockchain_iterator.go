package main 

import (
	"log"
	"github.com/boltdb/bolt"
)
/*
创建一个区块链的迭代器
*/
type BlockchainIterator struct{
	currentHash []byte 
	db *bolt.DB
}

//区块链迭代
func (bci *BlockchainIterator) Next() *Block{
	var block *Block
	err := bci.db.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte(blockTable))
		blockBytes := bucket.Get( bci.currentHash )
		block = UnSerializeBlock(blockBytes)
		return nil 
	})
	if err != nil {
		log.Panic(err)
	}
	//更新当前迭代器
	bci.currentHash = block.PrevBlockHash
	return block
}
