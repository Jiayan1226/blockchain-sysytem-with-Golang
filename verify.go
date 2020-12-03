package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
)

/*
交易模型：
1.拷贝修剪的副本
2.遍历原始交易，注意不是txCopy
3.找到output的公钥哈希，赋值给txCopy对应的input
4.还原签名的数据
//清理工作
5.校验
//还原签名r,s
*/
func (tx *Transaction)Verify() bool {
	fmt.Printf("对交易进行校验...\n")
	//1.拷贝修剪的副本
	txCopy:=tx.TrimmedCopy()
	//2.遍历原始交易，注意不是txCopy

	//3.找到output的公钥哈希，赋值给txCopy对应的input
	//4.还原签名的数据
	txCopy.SetTXID()
	//清理工作,
    txCopy.Txinput.PubKey=nil
	verifyData:=txCopy.Txid
	fmt.Printf("verifyData : %x\n",verifyData)
	//5.校验
	//还原签名r,s
	input:=tx.Txinput
	signature:=input.Isignature
	fmt.Printf("signature : %x\n",signature)
	fmt.Printf("Isignature长度为 : %d\n",len(signature))
	r:=big.Int{}
	s:=big.Int{}

	rData:=signature[:len(signature)/2]
	sData:=signature[len(signature)/2:]

	r.SetBytes(rData)
	s.SetBytes(sData)

	//还原公钥 curv,x,y
	//公钥字节流

	pubKeyBytes:=input.PubKey
	fmt.Printf("pubKeyBytes长度为 : %d",len(pubKeyBytes))
	x:=big.Int{}
	y:=big.Int{}

	xData:=pubKeyBytes[:len(pubKeyBytes)/2]
	yData:=pubKeyBytes[len(pubKeyBytes)/2:]

	x.SetBytes(xData)
	y.SetBytes(yData)

	curve:=elliptic.P256()

	publicKey:=ecdsa.PublicKey{curve,&x,&y}

	//数据，签名，公钥准备完毕，开始校验
	if !ecdsa.Verify(&publicKey,verifyData,&r,&s){
		return false
	}

	return true
}