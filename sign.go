package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
)

//非对称加密算法签名/验证包括三步：
//1.密钥生成Keygen （公钥、私钥）
//2.签名sign(r,s) 根据计算公式
//3.验证verify

//ECDSA 需要根据曲线和随机说来找到公钥和私钥
//ECDSA256表示ECDSA在签名时的算法是sha256，
func (tx *Transaction)Sign(privKey *ecdsa.PrivateKey) []byte {
	fmt.Printf("对交易进行签名...\n")
	//1.拷贝一份交易，txCopy,作相应裁剪，
	txCopy:=tx.TrimmedCopy()
	//2.遍历txCopy.inputs,签名要对数据的Hash进行签名，数据在交易中，要求交易的hash，就是Txid，所以我们可以使用交易id作为我们签名的内容
	//3.生成要签名的数据(哈希)
	txCopy.SetTXID()
	signData:=txCopy.Txid
	txCopy.Txinput.PubKey= nil
	fmt.Printf("要签名的数据,signData:%x\n",signData)
	//4.对数据进行签名r,s
	r,s,err:=ecdsa.Sign(rand.Reader,privKey,signData)
	if err!=nil{
		fmt.Printf("交易签名失败，err:%v\n",err)
	}
	//5.拼接r,s为字节流，赋值给原始交易的Signature字段
	signature:=append(r.Bytes(),s.Bytes()...)
	tx.Txinput.Isignature = signature

	fmt.Printf("签名结果 ： %x\n",signature)
	return signature

}

func (tx *Transaction)TrimmedCopy() Transaction {
	var input TxInput
	input = TxInput{tx.Txinput.Icoupon,tx.Txinput.IValue,nil,nil}
	output:=tx.Txoutput
	tx1:=Transaction{nil,input,output}
	return tx1
}