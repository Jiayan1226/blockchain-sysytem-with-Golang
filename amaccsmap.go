package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

//
type Maccs struct {
	Maccsmap  map[int64]*MerchantAccount
}

//创建Maccs，返回Maccs实例
func NewMaccs() *Maccs {
	var maccs Maccs
	maccs.Maccsmap = make(map[int64]*MerchantAccount)
	if !maccs.LoadFile(){
		fmt.Printf("加载商户钱包数据失败!\n")
	}
	return &maccs
}

const MerFile = "merchant.dat"
func (maccs *Maccs)SaveMer() bool {
	var buffer bytes.Buffer
	//将接口类型明确注册一下，否则gob编码失败！
	gob.Register(elliptic.P256())
	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(maccs)
	if err!=nil{
		fmt.Printf("商户序列化失败！\n")
		return false
	}
	content:=buffer.Bytes()
	err = ioutil.WriteFile(MerFile,content,0600)
	if err!=nil{
		fmt.Printf("商户账户序列化失败！\n")
		return false
	}
	return true
}

//把旧文件读出来，再写入数据
func (maccs *Maccs)LoadFile() bool {

	if !IsFileExist(MerFile){
		fmt.Printf("钱包文件不存在，准备创建!\n")
		return true
	}
	content,err:= ioutil.ReadFile(MerFile)
	if err!=nil{
		return false
	}
	gob.Register(elliptic.P256())
	decoder:=gob.NewDecoder(bytes.NewReader(content))
	var mas Maccs

	err = decoder.Decode(&mas)
	if err!=nil{
		fmt.Printf("err:%v\n",err)
		return false
	}
	maccs.Maccsmap = mas.Maccsmap

	return true
}

func (maccs *Maccs)Create() int64 {
	macc:= NewMerchantAccount1()
	address := GetRandom()
	maccs.Maccsmap[address] = macc
	res:=maccs.SaveMer()
	if !res{
		fmt.Printf("创建商户账户失败!\n")
		return 0
	}
	return address
}

func (maccs *Maccs)ListMaccs() []int64  {
	var addresses []int64
	for address,_:= range maccs.Maccsmap{
		addresses = append(addresses,address)
	}
	return addresses
}