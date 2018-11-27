package main 

import (
	"bytes"
	"crypto/sha256"
	"log"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"golang.org/x/crypto/ripemd160"
)
const addressCheckSum = 4
var version = byte(0x00)
//钱包结构体包含公私钥对
type Wallet struct{
	//私钥
	PrivateKey ecdsa.PrivateKey
	//公钥(该属性可以从私钥中获取)
	PublicKey []byte 
}
//构造钱包
func NewWallet() *Wallet{
	private,public := newPrivAndPub()
	return &Wallet{private,public}
}
//生成公私钥对
func newPrivAndPub()(ecdsa.PrivateKey,[]byte){
	//声明p-256曲线
	p256curve := elliptic.P256()
	//使用p256曲线生成私钥
	privKey,err := ecdsa.GenerateKey(p256curve,rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//在私钥获取公钥
	pubKey := append(privKey.PublicKey.X.Bytes(),
					 privKey.PublicKey.Y.Bytes()...)
	return *privKey,pubKey
}
//获取地址
func (wallet *Wallet) GetAddress() []byte{
	//对公钥进行Ripemd160
	pubkey160 := HashPubKey(wallet.PublicKey)
	//组合version + pubkey160
	payload := append([]byte{version},pubkey160...)
	//取其前4个字节
	checkSum := checksum(payload)
	//组合version + pubkey160 + checksum
	fullPayload := append(payload,checkSum...)
	//进行base58编码
	return Base58Encode(fullPayload)
}
//对公钥进行Ripemd160
func HashPubKey(pubkey []byte) []byte{
	pubkey256 := sha256.Sum256(pubkey)
	ripemd160Hasher := ripemd160.New()
	_,err := ripemd160Hasher.Write(pubkey256[:])
	if err != nil {
		log.Panic(err)
	}
	pubkey160 := ripemd160Hasher.Sum(nil)
	return pubkey160
}
//把version + pubkey160的结果进行两次sha256后取其前4个字节
func checksum(payload []byte) []byte{
	first256 := sha256.Sum256(payload)
	sec256 := sha256.Sum256(first256[:])
	return sec256[:addressCheckSum]
}

//验证地址是否正确
func ValidateAddress(address []byte) bool {
	pubKeyHash := Base58Decode(address)
	srcCheckSum := pubKeyHash[len(pubKeyHash) - addressCheckSum:]
	version := pubKeyHash[0]
	pubKey160 := pubKeyHash[1:len(pubKeyHash) - addressCheckSum]
	checkSum := checksum( append([]byte{version},pubKey160...) )
	return bytes.Compare(srcCheckSum,checkSum) == 0
}
