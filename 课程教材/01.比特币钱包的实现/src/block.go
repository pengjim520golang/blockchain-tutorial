package main

//区块数据结构
type Block struct {
	TimeStamp int64
	Hash []byte 
	PrevHash []byte 
	Nonce int 
	Transactions []*Transaction
}

//