package main

import (
	"log"
	"encoding/gob"
	"bytes"
)

type TXOutputs struct{
	Outputs []TXOutput
}
//序列化交易输出集合
func (outs TXOutputs) Serialize() []byte{
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
//反序列化交易输出集合
func UnSerializeTXOutputs(txOutputsBytes []byte) TXOutputs{
	var outputs TXOutputs 
	decoder := gob.NewDecoder( bytes.NewReader(txOutputsBytes) )
	err := decoder.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return outputs
}