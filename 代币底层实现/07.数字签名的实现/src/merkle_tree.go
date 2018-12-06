package main

import (
	"crypto/sha256"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left *MerkleNode
	Right *MerkleNode
	Data []byte
}

//创建默克尔树节点
func NewMerkleNode(left,right *MerkleNode,data []byte) *MerkleNode{
	merkleNode := &MerkleNode{}
	//如果没有左右节点,代表为叶子节点
	if left==nil && right==nil {
		hash := sha256.Sum256(data)
		merkleNode.Data = hash[:]
	}else{
	//如果拥有了左右节点，则代表叶子节点已经形成
		prevHashes := append(left.Data,right.Data...)
		hash := sha256.Sum256(prevHashes)
		merkleNode.Data = hash[:]
	}
	
	merkleNode.Left = left 
	merkleNode.Right = right

	return merkleNode
}

//创建默克尔树
func NewMerkleTree(data [][]byte) *MerkleTree{

	var nodes []MerkleNode
	//如果交易不是偶数则复制最后一个元素
	if len(data) % 2 != 0 {
		data = append( data , data[ len(data)-1 ] )
	}

	//组合叶子节点
	for _,datum := range data {
		node :=  NewMerkleNode(nil,nil,datum)
		nodes = append(nodes,*node)
	}

	for i:=0 ; i<len(data)/2 ; i++ {
		var newLevel []MerkleNode

		for j:=0;j<len(nodes);j+=2{
			node := NewMerkleNode(&nodes[j],&nodes[j+1],nil)
			newLevel = append(newLevel,*node)
		}

		nodes = newLevel
	}

	merkleTree := &MerkleTree{&nodes[0]}

	return merkleTree
}

