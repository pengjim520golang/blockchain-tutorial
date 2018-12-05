package main

import (
	"bytes"
)

type TXOutput struct {
	Value int
	PubKeyHash []byte
}

//创建输出
func NewTXOutput(value int , address string) *TXOutput{
	out := &TXOutput{value,nil}
	out.Lock([]byte(address))
	return out 
}

//锁定地址
func (out *TXOutput) Lock(address []byte){
	pubKey160 := Base58Decode(address)
	pubKey160 = pubKey160[1:len(pubKey160)-4]
	out.PubKeyHash = pubKey160
}

//解锁地址
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash,pubKeyHash) == 0
}

