package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
	"time"
)

type TxInput struct {
	Icoupon *CouPon
	IValue float64
	//IMid uint64
	PubKey []byte
	Isignature []byte

}

type TxOutput struct {
     OValue []byte
     //OMid uint64
     PubKeyHash []byte
     Ocoupon *CouPon
     time uint64
}

type Transaction struct {
	Txid []byte
	Txinput TxInput
	Txoutput TxOutput
}


//设置交易id
func (tx *Transaction)SetTXID()  {
	var buffer bytes.Buffer
	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	//交易内容有可能相同，加上时间戳区分

	hash:=sha256.Sum256(buffer.Bytes())
	tx.Txid=hash[:]
}
//根据地址找到公钥Hash
func (output *TxOutput)Lock(macc *MerchantAccount)  {
	output.PubKeyHash = HashPubKey(macc.Mwallets.PublicKey)
}

/*
	s = value
	1.首先判断是否有优惠券
	   如果有
	        //有可能优惠券已经花费但是金额不够
	         s - 优惠券金额
	         优惠券添加到商户中
	         同时优惠券置零：4种数据
	    如果没有
	        s
	2.判断顾客余额和s
	   如果付款金额>s
	        付款金额 - s
	        收款人金额 + s
	   否则
	       提示交易失败
	3.设置交易输入输出，并创建交易
	4.设置交易id
	5.返回交易结构
*/
//付款方，收款方，付款金额，优惠券的索引值
func NewCoinbase(clientacc *ClientAccount,merchantacc *MerchantAccount,value float64,i uint64) *Transaction  {
	s := value
	//1.首先判断是否有优惠券
	if clientacc.Ccoupons[i].CPid!=nil{
		s=s-clientacc.Ccoupons[i].Value
		merchantacc.ToCouPon(clientacc.Ccoupons[i])
		clientacc.DelCouPon(i)
	}
	//   如果有
	//        //有可能优惠券已经花费但是金额不够
	//         s - 优惠券金额
	//         优惠券添加到商户中
	//         同时优惠券置零：4种数据
	//    如果没有
	//        s
	//2.判断顾客余额和s
	if clientacc.CUtxo>=s{
		clientacc.CUtxo = clientacc.CUtxo - s
		merchantacc.MUtxo = merchantacc.MUtxo + s
	}else{
		fmt.Printf("交易失败！")
	}
	//   如果付款金额>s
	//        付款金额 - s
	//        收款人金额 + s
	//   否则
	//       提示交易失败
	//3.设置交易输入输出，并创建交易
	            //将交易设置成优惠券就可用
	input:= TxInput{clientacc.Ccoupons[i],s,merchantacc.Mwallets.PublicKey,nil}
    //将s转成字符串
   string:=FtoString(s)
   crypted,err:=AesEncrypt([]byte(string),merchantacc.Maeskey)
   if err!=nil{
   	fmt.Println(err)
   	fmt.Printf("加密成功！")
   }
	output:=TxOutput{crypted,nil,nil,uint64(time.Now().Unix())}
	output.Lock(merchantacc)
	tx:=Transaction{nil,input,output}
	//4.设置交易id
	tx.SetTXID()
	//交易签名,填写交易
	tx.Sign(clientacc.Cwallet.PrivateKey)
	//5.返回交易结构
	return &tx
}
func NewTransaction(from string,to string,value float64,i uint64,bc *Blockchain) *Transaction  {
	s := value

	maccs:= NewMaccs()
	merchantacc:= maccs.Maccsmap[StoInt64(from)]
	//macc.GetMacc()
	caccs:=NewCaccs()
	clientacc:=caccs.Caccsmap[StoInt64(to)]
	//1.首先判断是否有优惠券
	if clientacc.Ccoupons[i].CPid!=nil{
		s=s-clientacc.Ccoupons[i].Value
		merchantacc.ToCouPon(clientacc.Ccoupons[i])
		clientacc.DelCouPon(i)
	}
	//   如果有
	//        //有可能优惠券已经花费但是金额不够
	//         s - 优惠券金额
	//         优惠券添加到商户中
	//         同时优惠券置零：4种数据
	//    如果没有
	//        s
	//2.判断顾客余额和s
	if clientacc.CUtxo>=s{
		clientacc.CUtxo = clientacc.CUtxo - s
		merchantacc.MUtxo = merchantacc.MUtxo + s
		fmt.Printf("交易成功！")
	}else{
		fmt.Printf("交易失败！")
	}
	//   如果付款金额>s
	//        付款金额 - s
	//        收款人金额 + s
	//   否则
	//       提示交易失败
	//3.设置交易输入输出，并创建交易
	input:= TxInput{clientacc.Ccoupons[i],s,clientacc.Cwallet.PublicKey,nil}
	//将s转成字符串
	string:=FtoString(s)
	crypted,err:=AesEncrypt([]byte(string),merchantacc.Maeskey)
	if err!=nil{
		fmt.Println(err)
		fmt.Printf("加密成功！")
	}
	fmt.Printf("测试：创建新的交易！")
	output:=TxOutput{crypted,nil,nil,uint64(time.Now().Unix())}
	output.Lock(merchantacc)
	tx:=Transaction{nil,input,output}
	//4.设置交易id
	tx.SetTXID()

	//交易签名,填写交易
	mon:=MatchNode(from)
	mon.SignTransaction(&tx,clientacc.Cwallet.PrivateKey)
	//5.返回交易结构
	return &tx
}
func (tx *Transaction)String() string  {
	var lines []string

	lines = append(lines,fmt.Sprintf("---Transaction %x : ",tx.Txid))
	lines = append(lines,fmt.Sprintf("    Input:"))
    lines = append(lines,fmt.Sprintf("    Value  :         %f ",tx.Txinput.IValue))
    lines = append(lines,fmt.Sprintf("    PubKey :         %x",tx.Txinput.PubKey))
	lines = append(lines,fmt.Sprintf("    Signature:       %x",tx.Txinput.Isignature))
	if tx.Txinput.Icoupon!=nil  && tx.Txinput.Icoupon.CPid!=nil {
		lines = append(lines, fmt.Sprintf("    Coupon.Cpid:     %x", tx.Txinput.Icoupon.CPid))
		lines = append(lines, fmt.Sprintf("    Coupon.Value:    %f", tx.Txinput.Icoupon.Value))
	}else {
		fmt.Printf("nil\n")
	}
	if tx.Txinput.Icoupon.Gp!=nil {
		lines = append(lines, fmt.Sprintf("    Coupon.Gp.Mid:   %d", tx.Txinput.Icoupon.Gp.Mid))
		lines = append(lines, fmt.Sprintf("    Coupon.Gp.Gmt:   %d", tx.Txinput.Icoupon.Gp.Gmt))
	}else{
		fmt.Printf("nil\n")
	}


	lines = append(lines,fmt.Sprintf("     Onput :"))
	lines = append(lines,fmt.Sprintf("     Value :            %x",tx.Txoutput.OValue))
	lines = append(lines,fmt.Sprintf("     PubKeyHash :       %x",tx.Txoutput.PubKeyHash))
	if tx.Txinput.Icoupon!=nil  && tx.Txinput.Icoupon.CPid!=nil {
		lines = append(lines, fmt.Sprintf("     Coupon.Cpid :      %x", tx.Txoutput.Ocoupon.CPid))
		lines = append(lines, fmt.Sprintf("     Coupon.Value :     %f", tx.Txoutput.Ocoupon.Value))
	}else {
		fmt.Printf("nil\n")
	}
	if tx.Txinput.Icoupon.Gp!=nil {
		lines = append(lines, fmt.Sprintf("     Coupon.Gp.Mid :    %d", tx.Txoutput.Ocoupon.Gp.Mid))
		lines = append(lines, fmt.Sprintf("     Coupon.Gp.Gmt :    %d", tx.Txoutput.Ocoupon.Gp.Gmt))
	}else{
		fmt.Printf("nil\n")
	}


	return strings.Join(lines,"\n")

}
