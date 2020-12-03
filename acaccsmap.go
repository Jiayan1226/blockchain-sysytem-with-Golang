package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

//顾客集合
type Caccs struct {
	Caccsmap  map[int64]*ClientAccount
}

//创建Maccs，返回Maccs实例
func NewCaccs() *Caccs {
	var caccs Caccs
	caccs.Caccsmap = make(map[int64]*ClientAccount)
	if !caccs.LoadcFile(){
		fmt.Printf("加载顾客钱包数据失败!\n")
	}
	return &caccs
}

const CerFile = "Cerchant.dat"
func (caccs *Caccs)SaveCer() bool {
	var buffer bytes.Buffer
	//将接口类型明确注册一下，否则gob编码失败！
	gob.Register(elliptic.P256())
	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(caccs)
	if err!=nil{
		fmt.Printf("顾客序列化失败！\n	")
		return false
	}
	content:=buffer.Bytes()
	err = ioutil.WriteFile(CerFile,content,0600)
	if err!=nil{
		fmt.Printf("顾客账户序列化失败！\n")
		return false
	}
	return true
}

//把旧文件读出来，再写入数据
func (caccs *Caccs)LoadcFile() bool {

	if !IsFileExist(CerFile){
		fmt.Printf("钱包文件不存在，准备创建!\n")
		return true
	}
	//读取文件
	content,err:= ioutil.ReadFile(CerFile)
	if err!=nil{
		return false
	}
	gob.Register(elliptic.P256())
	decoder:=gob.NewDecoder(bytes.NewReader(content))
	var cas Caccs

	err = decoder.Decode(&cas)
	if err!=nil{
		fmt.Printf("err:%v\n",err)
		return false
	}
	caccs.Caccsmap = cas.Caccsmap

	return true
}

func (caccs *Caccs)Create() int64 {
	cacc:= NewClientAccount()
	address := GetRandom()
	caccs.Caccsmap[address] = cacc
	res:=caccs.SaveCer()
	if !res{
		fmt.Printf("创建顾客账户失败!\n")
		return 0
	}
	return address
}

func (caccs *Caccs)ListCaccs() []int64  {
	var addresses []int64
	for address,_:= range caccs.Caccsmap{
		addresses = append(addresses,address)
	}
	return addresses
}