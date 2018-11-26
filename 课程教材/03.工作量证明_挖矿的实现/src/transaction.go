package main

import (
	"log"
	"encoding/gob"
	"bytes"
)

type Transaction struct{
	ID []byte
	Vint []TXInput
	Vout []TXOutput
}

func (tx Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}