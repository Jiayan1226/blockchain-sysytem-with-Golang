package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type WalletKeyPair struct {
	PrivateKey *ecdsa.PrivateKey  //私钥
	//我们可以将公钥X,Y进行字节流拼接后进行传输，这样对端进行切割还原，
	PublicKey []byte //公钥中包含曲线，和某一点B的坐标X和Y，如果是SHA256 那么X，Y分别是256位。
}

func NewWalletKeyPair() *WalletKeyPair  {
	privateKey,err:=ecdsa.GenerateKey(elliptic.P256(),rand.Reader)//根据曲线随机数生成，私钥中包括公钥，以及随机数D，公钥中包括，根据D生成公钥中x,y值，以及椭圆曲线。
	if err!=nil{
		log.Panic(err)
	}
	publicKeyRaw:=privateKey.PublicKey
	publicKey:=append(publicKeyRaw.X.Bytes(),publicKeyRaw.Y.Bytes()...)
	return &WalletKeyPair{PrivateKey:privateKey,PublicKey:publicKey}
}

//获取公钥哈希
func HashPubKey(pubKey []byte) []byte {
	hash:=sha256.Sum256(pubKey)
	//创建一个hash160对象
	//向Hash160中write数据
	//做哈希运算
	rip160Hasher:=ripemd160.New()
	_,err:=rip160Hasher.Write(hash[:])

	if err!=nil{
		log.Panic(err)
	}
	//Sum函数会把我们的结果与Sum参数append到一起，然后返回，我们传入nil，防止数据污染
	publicHash:=rip160Hasher.Sum(nil)
	return publicHash
}