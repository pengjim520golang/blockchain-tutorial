package main 

import (
	"fmt"
	"crypto/sha256"
	"math"
	"bytes"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

//挖矿的准则：求一个比目标值小大的哈希值
//targetBit越小代表目标值越大,目标值越大则挖矿速度越快
//targetBit越大代表目标值越小,目标值越小则挖矿速度越慢
const targetBits = 20

type ProofOfWork struct {
	block *Block
	target *big.Int
}

//创建一个工作量证明的实例
func NewProofOfWork(block *Block) *ProofOfWork{
	//创建一个256长度大数
	target := big.NewInt(1)
	//如果targetBit为0时代表target为最大值,如果targetBit为256时为最小值
	target.Lsh(target,uint(256-targetBits))
	//创建工作量证明实例
	pow := &ProofOfWork{block,target}

	return pow 
}

//创建用于比较目标值的哈希数据
func (pow *ProofOfWork) prepareData(nonce int) []byte{
	data := bytes.Join([][]byte{pow.block.PrevBlockHash,pow.block.HashTransactions(),IntToHex(pow.block.Timestamp),IntToHex(int64(targetBits)),IntToHex(int64(nonce))},[]byte{})
	return data 
}

//挖矿
func (pow *ProofOfWork) Run()(int,[]byte){
	var hashInt big.Int
	var hash [32]byte 
	nonce := 0
	fmt.Printf("正在挖矿中...")
	for nonce<maxNonce{
		//通过不断改变nonce的值来求出不同的哈希
		data := pow.prepareData(nonce)
		//把数据变为哈希
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		//将哈希转为可以比较的大整数
		hashInt.SetBytes(hash[:])
		// Cmp返回-1表示比目标值小
		if hashInt.Cmp(pow.target) == -1{
			//如果比目标值小代表挖矿成功
			break 
		}else{
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce,hash[:]
}
//验证工作量证明是否正确
func (pow *ProofOfWork) Validate() bool{
	var hashInt big.Int
	var hash [32]byte
	data := pow.prepareData( pow.block.Nonce )
	hash = sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1 
}