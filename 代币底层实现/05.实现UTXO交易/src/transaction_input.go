package main 

import (
	"bytes"
)

//交易输入
type TXInput struct{
	Txid []byte 
	Vout int 
	Signature []byte
	PubKey []byte 
}

//消费一笔输出时进行检查
func (in *TXInput) UsesKey(PubKeyHash []byte) bool{
	PubKeyHash160 := HashPubKey(in.PubKey)
	return bytes.Compare(PubKeyHash160,PubKeyHash) == 0
}