package main

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Version uint64
	MerKleRoot []byte
	Timestamp uint64
	Transactions []*Transaction
    PrevBlockHash []byte
	Hash []byte
}

func NewBlock(txs []*Transaction,prevblockhash []byte) *Block {
	block:= Block{
		Version:00,
		MerKleRoot:[]byte{},
		Timestamp:uint64(time.Now().Unix()),
		Transactions:txs,
		PrevBlockHash:prevblockhash,
	}
	block.HashTransactions()
    block.PrePareHash()

	return &block
}

func (block *Block)HashTransactions()  {
	var hashes []byte
	for _,tx:=range block.Transactions{
		txid/*[]byte*/:=tx.Txid
		hashes=append(hashes,txid...)
	}
	hash:=sha256.Sum256(hashes)
	block.MerKleRoot=hash[:]
}
func (block *Block)PrePareHash()  {
	data:= bytes.Join([][]byte{
		block.MerKleRoot,
		intToByte(block.Timestamp),
		block.PrevBlockHash,
	},[]byte{})
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}

