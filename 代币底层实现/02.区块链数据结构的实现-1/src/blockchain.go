package main 

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(transactions []*Transaction){
	//获取上一个
	prevBlock := bc.blocks[ len(bc.blocks) - 1 ]
	//创建一个区块
	newBlock := NewBlock(transactions,prevBlock.Hash)
	//把区块加入区块链中
	bc.blocks = append(bc.blocks,newBlock)
}

//创建拥有创世区块的区块链
func NewBlockchain(block *Block) *Blockchain{
	return &Blockchain{[]*Block{block}}
}