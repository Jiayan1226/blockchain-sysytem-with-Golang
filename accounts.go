package main

import (
	"fmt"
)

type GuidePoint struct {
	Gmt uint64
	Mid uint64
}

type CouPon struct {
	CPid []byte
	Value float64
	Gp *GuidePoint
}

type ClientAccount struct {
	 Cid  uint64
     CUtxo float64
	 Cwallet *WalletKeyPair
     Ccoupons []*CouPon
}

type MerchantAccount struct {
	 Mid  uint64
     MUtxo float64
	 Maeskey []byte
	 Mwallets *WalletKeyPair
     Mcoupons []*CouPon
     Mguidepoints []*GuidePoint
}

type ShoppingMallAccount struct {
     SUtxo float64
     SGuidepoints []*GuidePoint
}

//初始化一个积分
func NewGuidePoint()  *GuidePoint{
	guidepoint:=GuidePoint{50,2001}
	return &guidepoint
}

func (sacc *ShoppingMallAccount)NewGuidePoint(mer *MerchantAccount) {
	guidepoint:=GuidePoint{50,mer.Mid}
	mer.Mguidepoints= append(mer.Mguidepoints,&guidepoint)
}

//初始化一个优惠券
func NewCouPon() *CouPon{
	gp:=NewGuidePoint()
	coupon:=CouPon{[]byte("Cp01"),5,gp}
	return &coupon
}

func (macc *MerchantAccount)NewCouPon(cacc *ClientAccount)  {

}
//初始化一个顾客
func NewClientAccount() *ClientAccount {
	coupon:=NewCouPon()
	coupons:=[]*CouPon{coupon}
	ws:=NewWalletKeyPair()
	clientaccount:=ClientAccount{
		Cid:1001,
		CUtxo:5000,
		Cwallet:ws,
		Ccoupons:coupons,
	}
	return &clientaccount

}
//初始化一个商户
func NewMerchantAccount1() *MerchantAccount {
	mcoupons:=[]*CouPon{NewCouPon()}
	mgps:=[]*GuidePoint{}
	ws:=NewWalletKeyPair()
    merchantaccount:=MerchantAccount{
   	Mid:2001,
   	MUtxo:2000,
   	Maeskey:[]byte("vdncloud12345678"),
   	Mwallets:ws,
   	Mcoupons:mcoupons,
   	Mguidepoints:mgps,
	}
   return &merchantaccount
}
func NewMerchantAccount2() *MerchantAccount {
	mcoupons:=[]*CouPon{NewCouPon()}
	mgps:=[]*GuidePoint{}
	merchantaccount:=MerchantAccount{
		Mid:2002,
		MUtxo:3000,
		Maeskey:[]byte("vdncloud123456"),
		Mcoupons:mcoupons,
		Mguidepoints:mgps,
	}
	return &merchantaccount
}
//初始化一个购物中心
func NewShoppingMallAccount()  *ShoppingMallAccount{
	sguidepoints:=[]*GuidePoint{}
	shoppingmallaccount:=ShoppingMallAccount{
		SUtxo:3000,
		SGuidepoints:sguidepoints,
	}
	return &shoppingmallaccount
}
//在商户中添加优惠券
func (macc *MerchantAccount)ToCouPon(coupon *CouPon)  {
	macc.Mcoupons = append(macc.Mcoupons, coupon)
}
//在顾客中添加优惠券
func (cacc *ClientAccount)ToCouPon(coupon *CouPon)  {
	cacc.Ccoupons = append(cacc.Ccoupons, coupon)
}
//在顾客中剔除优惠券
func (cacc *ClientAccount)DelCouPon(i uint64)  {
	coupon:=CouPon{}
	cacc.Ccoupons[i]=&coupon
}

//获取顾客账户
func (cacc *ClientAccount)GetCacc()  {
	fmt.Printf("Cid:%d     ",cacc.Cid)
	fmt.Printf("Utxo:%f    ",cacc.CUtxo)
	fmt.Printf("\n")
	fmt.Printf("PublicKey:%x    ",cacc.Cwallet.PublicKey)
	fmt.Printf("\n")
    if cacc.Ccoupons!=nil && cacc.Ccoupons[0].CPid!=nil{
		for _,coupon:= range cacc.Ccoupons{
			fmt.Printf("{ Cpid:%s    ",coupon.CPid)
			fmt.Printf("Value:%f   ",coupon.Value)
			fmt.Printf("Gp.id:%d   ",coupon.Gp.Mid)
			fmt.Printf("Gp.Gmt:%d } ",coupon.Gp.Gmt)
			fmt.Printf("\n")
		}
	}else{
		fmt.Printf("\n")
	}
}
//获取商户账户
func (macc *MerchantAccount)GetMacc()  {
	fmt.Printf("Mid:%d    ",macc.Mid)
	fmt.Printf("MUtxo:%f    ",macc.MUtxo)
	fmt.Printf("Maeskey:%s   ",macc.Maeskey)
	fmt.Printf("\n")
	fmt.Printf("PublicKey:%x   ",macc.Mwallets.PublicKey)
	fmt.Printf("\n")
	if macc.Mcoupons!=nil  && macc.Mcoupons[0].CPid!=nil{
		for _,coupon:= range macc.Mcoupons {
			fmt.Printf("{ Cpid:%s    ",coupon.CPid)
			fmt.Printf("Value:%f   ",coupon.Value)
			fmt.Printf("Gp.id:%d   ",coupon.Gp.Mid)
			fmt.Printf("Gp.Gmt:%d } ",coupon.Gp.Gmt)
			fmt.Printf("\n")
		}
	}else {
		fmt.Printf("nil\n")
	}
    if macc.Mguidepoints!=nil{
		for _,gp:=range macc.Mguidepoints{
			fmt.Printf("{ Mid:%d    ",gp.Mid)
			fmt.Printf("{ Gmt:%d    ",gp.Gmt)
			fmt.Printf("\n")
		}
	}else{
		fmt.Printf("nil\n")
	}

}
//获取购物中心账户
func GetSacc()  {

}
