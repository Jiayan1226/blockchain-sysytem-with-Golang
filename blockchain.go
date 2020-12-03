package main

import (
	"fmt"
	 "github.com/boltdb/bolt"
	"log"
	"os"
)

type Blockchain struct {
	db *DB
	tail []byte
}

const blockChainName  = "blockChain.db"
const blockBucketName  = "blockBucket"
const lastHashKey  = "lastHashKey"

func (bc *Blockchain)AddBlock(con *ConsensusNode,txs []*Transaction)  {
	//验证
    validTXs:=[]*Transaction{}
    for _,tx:=range txs{
    	if con.VerifyTransaction(tx){
    		fmt.Printf("该交易有效 : %x\n",tx.Txid)
    		validTXs = append(validTXs,tx)
		}else{
			fmt.Printf("发现无效的交易: %x\n",tx.Txid)
		}
	}
	bc.db.Update(func(tx *Tx) error {
			b:=tx.Bucket([]byte(blockBucketName))//判断数据库有没有被调用
			if b==nil{
				fmt.Printf("bucket不存在，请检查!\n")
				os.Exit(1)
			}
			//抽屉准备完毕，添加创世区块，只有一个交易
			block:=NewBlock(txs,bc.tail)
			b.Put(block.Hash,block.Serialize())
			b.Put([]byte("lastHashKey"),block.Hash)
			bc.tail = block.Hash
		return nil

	})
}
//返回区块链实例，保证区块链已经创建完毕

func NewBlockchain()*Blockchain  {
	if !IsFileExist(blockChainName){
		fmt.Printf("区块链不存在，请先创建！\n")
		return nil
	}
	//获取数据库的句柄，打开数据库，读写数据
	db,err:= Open(blockChainName,0600,nil)
	if err!=nil{
		log.Panic(err)
	}
	var tail []byte
	db.View(func(tx *Tx) error {
			b:=tx.Bucket([]byte("blockBucket"))
			if b==nil{
				fmt.Printf("区块链bucket为空，请检查！\n")
				os.Exit(1)
			}
			tail = b.Get([]byte("lastHashKey"))
		return nil
	})
	return &Blockchain{db,tail}
}
//首次创建区块链
func CreateBlockchain(clientacc *ClientAccount,merchantacc *MerchantAccount,value float64,i uint64) *Blockchain {
	if IsFileExist(blockChainName){
		fmt.Printf("区块链已经存在，不需要重复创建！\n")
		return nil
	}
	//获得数据库的句柄，打开数据库，读写数据
	db,err:= Open(blockChainName,0600,nil)
	if err!=nil{
		log.Panic(err)
	}
	var tail []byte
	db.Update(func(tx *Tx) error {
		  b,err:=tx.CreateBucket([]byte(blockBucketName))
		  if err!=nil{
		  	log.Panic(err)
		  }
		  //抽屉准备完毕，开始添加创世区块，创世区块中只有一个挖矿交易，只有coinbase
		  coinbase:=NewCoinbase(clientacc,merchantacc,value,i)
		  genesisBlock:=NewBlock([]*Transaction{coinbase},[]byte{})
		  b.Put(genesisBlock.Hash,genesisBlock.Serialize())
		  b.Put([]byte("lastHashKey"),genesisBlock.Hash)

		  tail = genesisBlock.Hash
		  return nil
	})

	return &Blockchain{db,tail}
}

