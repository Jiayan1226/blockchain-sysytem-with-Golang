package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

//生成指定长度的随机数
func CreateRandomNumber(len int) string {
	var numbers = []byte{1,2,3,4,5,6,7,8,9}
	var container string

	length:=bytes.NewReader(numbers).Len()//创建一个从numbers读取数据的reader

	for i:=1;i<=len;i++{
		random,err:=rand.Int(rand.Reader,big.NewInt(int64(length)))
		if err!=nil{
			os.Exit(1)
		}
		container += fmt.Sprintf("%d",numbers[random.Int64()])
	}

	return container
}
//生成随机长度的字符串
func CreateRandomString(len int)string  {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b:=bytes.NewBufferString(str)//NewBuffer使用s作为初始内容创建并初始化一个Buffer、本函数用于创建一个用于读取已存在数据的buffer
	//大多数情况下，new(Buffer)(或只是声明一个Buffer类型变量)就足以初始化一个Buffer了
    length:=b.Len()
    bigInt:= big.NewInt(int64(length))
    for i:=0;i < len;i++{
    	randomInt,_:=rand.Int(rand.Reader,bigInt)
    	container+=string(str[randomInt.Int64()])
	}
    return container
}


