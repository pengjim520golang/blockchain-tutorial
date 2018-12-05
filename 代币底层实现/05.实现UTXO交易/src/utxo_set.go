package main

import (
	"encoding/hex"
	"log"
	"github.com/boltdb/bolt"
)

const UtxoTable = "utxoTable"

type UTXOSet struct{
	Blockchain *Blockchain
}

//根据哈希找到可以花费的输出
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte,amount int)(int,map[string][]int){
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	err := u.Blockchain.db.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte(UtxoTable))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			txID := hex.EncodeToString(k)
			txOutputs := UnSerializeTXOutputs(v)
			for outIndex,out := range txOutputs.Outputs{
				if out.IsLockedWithKey(pubKeyHash) && accumulated<amount{
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID],outIndex)
				}
			}
		}
		return nil  
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated,unspentOutputs
}

//根据一个地址的公钥找到所有输出
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput{
	var outputs []TXOutput
	
	err := u.Blockchain.db.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte(UtxoTable))
		cursor := bucket.Cursor()
		for k,v := cursor.First();k!=nil;k,v = cursor.Next(){
			//获取输出集合
			txOutputs := UnSerializeTXOutputs(v)
			//遍历输出集合
			for _,out := range txOutputs.Outputs{
				if out.IsLockedWithKey(pubKeyHash) {
					outputs = append(outputs,out)
				}
			}
		}
		return nil 
	})
	if err != nil {
		log.Panic(err)
	}
	return outputs
}

//创建索引
func (u UTXOSet) Reindex(){
	db := u.Blockchain.db
	//删除UtxoTable
	err := db.Update(func(tx *bolt.Tx) error{
		err := tx.DeleteBucket([]byte(UtxoTable))
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}
		return nil 
	})
	//找出区块链中所有的UTXO
	utxo := u.Blockchain.FindUTXO()
	//重新创建UtxoTable进行录入
	err = db.Update(func(tx *bolt.Tx) error{
		bucket,err := tx.CreateBucket([]byte(UtxoTable))
		if err != nil {
			log.Panic(err)
		}
		//录入到UtxoTable中
		for txid,outputs := range utxo{
			txID,err := hex.DecodeString(txid)
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put(txID,outputs.Serialize())
			if err != nil {
				log.Panic(err)
			}

		}
		return nil 
	})
	if err != nil {
		log.Panic(err)
	}
}

func (u UTXOSet) Update(block *Block){
	db := u.Blockchain.db
	err := db.Update(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte(UtxoTable))
		for _,tx := range block.Transactions{
			if tx.IsCoinbase() == false {
				//如果不是coinbase交易，必定引用了某笔交易的输出
				for _,vin := range tx.Vint{
					updateOuts := TXOutputs{}
					//找到对应交易的输出集合
					outsBytes := bucket.Get(vin.Txid)
					//反序列化
					outs := UnSerializeTXOutputs(outsBytes)
					//找到尚未消费的输出
					for outIndex,out := range outs.Outputs{
						if outIndex != vin.Vout{
							updateOuts.Outputs = append(updateOuts.Outputs,out)
						}
					}
					if len( updateOuts.Outputs ) == 0 {
						err := bucket.Delete(vin.Txid)
						if err != nil {
							log.Panic(err)
						}
					}else{
						err := bucket.Put(vin.Txid,updateOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}

			}// if IsCoinbase ending
			
			//把新的交易输出录入
			newOutputs := TXOutputs{}
			for _,out := range tx.Vout{
				newOutputs.Outputs = append(newOutputs.Outputs,out)
			}
			err := bucket.Put(tx.ID,newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		} //for _,tx := range block.Transactions ending
		return nil 
	})
	if err != nil {
		log.Panic(err)
	}
}