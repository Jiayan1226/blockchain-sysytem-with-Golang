package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//这是一个工具函数文件
func intToByte(num uint64)[]byte{
	//使用binary.Write来进行编码
	var buffer bytes.Buffer
	//编码要进行错误检查，一定要做
	err:=binary.Write(&buffer,binary.BigEndian,num)
	if err !=nil{
		log.Panic(err)
	}
	return buffer.Bytes()
}
//判断文件是否存在
func IsFileExist(fileName string)bool  {
	//使用os.Stat来判断
	_,err:=os.Stat(fileName)

	if os.IsNotExist(err){
		return false
	}
	return true
}

func (block *Block)Serialize()[]byte  {
	var buffer bytes.Buffer
	//定义编码器
	encoder:=gob.NewEncoder(&buffer)
	//编码器对结构进行编码，一定要进行校验
	err:=encoder.Encode(block)
	if err!=nil{
		log.Panic(err)
	}
	return buffer.Bytes()
}

func Deserialize(data []byte)*Block  {
	//fmt.Printf("解码传入的数据：%x\n",data)
	var block Block
	//创建解码器
	decoder:=gob.NewDecoder(bytes.NewReader(data))
	err:=decoder.Decode(&block)
	if err!=nil{
		log.Panic(err)
	}
	return &block
}

func GetRandom() int64{
	r:=rand.New(rand.NewSource(time.Now().UnixNano()))
	x:=r.Int63n(1000000000)
	return x
}

//string 转到int64
func StoInt64(string string) int64 {
	int64,err:=strconv.ParseInt(string,10,64)
	if err!=nil{
		log.Panic(err)
		return 0
	}
	return int64
}

//string转成float64
func StoFloat64(v string) float64 {
	v1,err:=strconv.ParseFloat(v,64)
	if err!=nil{
		log.Panic(err)
	}
	return v1
}

//将float64转成string
func FtoString(s float64) string  {
	string:=strconv.FormatFloat(s,'E',-1,64)
	return string
}