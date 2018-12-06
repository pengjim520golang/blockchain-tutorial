package main

import (
	"math/big"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"crypto/sha256"
	"log"
	"encoding/gob"
	"bytes"
	"fmt"
)
//基币奖励
const subsidy = 10

type Transaction struct{
	ID []byte
	Vint []TXInput
	Vout []TXOutput
}
//序列化交易
func (tx Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
//把当前交易进行哈希
func (tx *Transaction) Hash() []byte{
	var hash [32]byte
	txCopy := *tx 
	txCopy.ID = []byte{}
	hash = sha256.Sum256(tx.Serialize()) 
	return hash[:]
}
//创建coinbase交易
func NewCoinbaseTX( to,data string ) *Transaction{
	if data =="" {
		//data = "coinbase transaction"
		randData := make([]byte,20)
		_,err :=  rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}
	//创建一笔输入
	txIn := TXInput{[]byte{},-1,nil,[]byte(data)}
	txOut := NewTXOutput(subsidy,to) 
	cbtx := Transaction{nil,[]TXInput{txIn},[]TXOutput{*txOut}}
	cbtx.ID = cbtx.Hash()
	return &cbtx 
}
//判断当前交易是否属于Coinbase
func (tx *Transaction) IsCoinbase() bool{
	return (len(tx.Vint)==1 && tx.Vint[0].Vout==-1 && len(tx.Vint[0].Txid)==0)
}
//创建一笔交易
func NewUTXOTransaction(from , to string,amount int,utxoSet *UTXOSet) *Transaction{
	var Txinputs []TXInput
	var Txoutputs []TXOutput
	//获取钱包中的PubKey
	wallets,err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc,outputs := utxoSet.FindSpendableOutputs(pubKeyHash,amount)
	if acc < amount {
		log.Panic("没有足够完成本次交易的金额")
	}
	for txid,outputIndexes := range outputs{
		txID,err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		//遍历可用的输出索引
		for _,outIndex := range outputIndexes{
			txIn := TXInput{txID,outIndex,nil,wallet.PublicKey}
			Txinputs = append(Txinputs,txIn)
		}
	}
	//创建一笔交易输出
	txOut := *NewTXOutput(amount,to)
	Txoutputs = append(Txoutputs,txOut)
	if acc > amount {
		change := acc - amount 
		//创建一逼找零输出
		changeOut := *NewTXOutput(change,from)
		Txoutputs = append(Txoutputs,changeOut)
	}
	transaction := Transaction{nil,Txinputs,Txoutputs}
	transaction.ID = transaction.Hash()
	//签名当前交易
	utxoSet.Blockchain.SignTransaction(&transaction,wallet.PrivateKey)
	return &transaction
}

//复制交易副本
func (tx *Transaction) TrimmedCopy() Transaction{
	var txOutputs []TXOutput
	var txInputs []TXInput
	//构造txInput集合
	for _,in := range tx.Vint{
		txIn := TXInput{in.Txid,in.Vout,nil,nil} 
		txInputs = append(txInputs,txIn)
	}
	//构造txOutput集合
	for _,out := range tx.Vout{
		txOut := TXOutput{out.Value,out.PubKeyHash} 
		txOutputs = append(txOutputs,txOut)
	}
	//构造一个新的Transaction对象
	txCopy := Transaction{tx.ID,txInputs,txOutputs}
	return txCopy
}

//数字签名
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey,preTxs map[string]Transaction){
	//coinbase交易无需签名
	if tx.IsCoinbase() {
		return 
	}

	//验证当前交易中所有的引用是否有效
	for _,in := range tx.Vint{
		if preTxs[hex.EncodeToString(in.Txid)].ID == nil {
			log.Panic("当前的交易引用不存在")
		}
	}

	//复制当前交易的副本
	txCopy := tx.TrimmedCopy()
	//对引用进行签名
	for inIndex,txIn := range txCopy.Vint {
		prevTx := preTxs[hex.EncodeToString(txIn.Txid)]
		txCopy.Vint[inIndex].Signature = nil 
		txCopy.Vint[inIndex].PubKey = prevTx.Vout[txIn.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vint[inIndex].PubKey = nil 

		r,s,err := ecdsa.Sign(rand.Reader,&privateKey,txCopy.ID) 
		if err != nil {
			log.Panic(err)
		}
		//拼接r,s就是数字签名
		Signature := append(r.Bytes(),s.Bytes()...)
		//fmt.Printf("签名：%x\n",Signature)
		tx.Vint[inIndex].Signature = Signature
	}
}

//验证数字签名
func (tx *Transaction) Verify(preTxs map[string]Transaction) bool{
	if tx.IsCoinbase(){
		return true 
	}

	for _,in := range tx.Vint{
		if preTxs[ hex.EncodeToString(in.Txid) ].ID == nil {
			log.Panic("无法找到引用交易")
		}
	}

	//复制当前交易
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inIndex,txIn := range tx.Vint{
		prevTx := preTxs[ hex.EncodeToString(txIn.Txid) ]
		txCopy.Vint[inIndex].Signature = nil 
		txCopy.Vint[inIndex].PubKey = prevTx.Vout[txIn.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vint[inIndex].PubKey = nil 

		//重新取出签名中的r和s
		r := big.Int{}
		s := big.Int{}
		signLen := len(txIn.Signature) 
		//fmt.Printf("in sign的%x\n",txIn.Signature)
		r.SetBytes(txIn.Signature[:(signLen/2)])
		s.SetBytes(txIn.Signature[(signLen/2):])

		//提取钱包中的公钥x和y
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(txIn.PubKey)
		//尝试构造和钱包中一样的PubKey
		x.SetBytes(txIn.PubKey[:(pubKeyLen/2)])
		y.SetBytes(txIn.PubKey[(pubKeyLen/2):])
		rawPubKey := ecdsa.PublicKey{curve,&x,&y}

		//fmt.Printf("in的%x\n",txIn.PubKey)
		//由于私钥构造了数字签名,公钥是从私钥中获取的,验证签名时就需要把公钥传入
		//公钥的验证是比对私钥中的PubKey是否一致
		//r,x是数字签名中提取的,因此也需要传入,txCopy也是构造签名的要素也需要传入
		//只要在交易中有一个引用是签名不同，整个交易就会回滚
		//fmt.Printf("txCopy的%x\n",txCopy.ID)
		if ecdsa.Verify(&rawPubKey,txCopy.ID,&r,&s) == false {
			return false
		}
	}
	return true 
}

