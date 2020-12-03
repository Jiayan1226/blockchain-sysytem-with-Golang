package main

import (
	"github.com/boltdb/bolt"
	"fmt"
)

type BlockChainIterator struct {
	db *bolt.DB
	current []byte //当前所指向的区块的哈希
}

func (bc *Blockchain)NewIterator() *BlockChainIterator  {
	return &BlockChainIterator{bc.db,bc.tail}
}

func (it *BlockChainIterator)Next() *Block  {
	var block Block
	it.db.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockBucketName))
		if b == nil{
			fmt.Printf("bucket不存在，请检查!\n")
		}
		//真正的读取数据
		blockInfo:=b.Get(it.current)
		block = *Deserialize(blockInfo)
		it.current = block.PrevBlockHash

		return nil
	})
	return &block
}