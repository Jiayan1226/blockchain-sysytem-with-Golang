package main

//const WalletName  = "wallet.dat"
//
////地址和密钥对应起来
//type Wallets struct {
//	WalletsMap map[string]*WalletKeyPair
//}
//
//func (ws *Wallets) LoadFromFile() bool {
//	//判断文件是否存在
//	if !IsFileExist(WalletName){
//		fmt.Printf("钱包文件不存在，准备创建！\n")
//		return true
//	}
//	content,err:=ioutil.ReadFile(WalletName)
//	if err!=nil{
//		return false
//	}
//
//	gob.Register(elliptic.P256())
//	//gob编码
//	decoder:=gob.NewDecoder(bytes.NewBuffer(content))
//	var wallets Wallets
//	err = decoder.Decode(&wallets)
//	if err!=nil{
//		fmt.Printf("err:%v\n",err)
//		return false
//	}
//	ws.WalletsMap = wallets.WalletsMap
//
//	return true
//}
//
////创建Wallets,返回Wallets的实例
//func NewWallets() *Wallets {
//	var ws Wallets
//	//1.把所有钱包从本地加载出来
//	ws.WalletsMap = make(map[string]*WalletKeyPair)
//	if !ws.LoadFromFile(){
//		fmt.Printf("加载钱包数据失败！\n")
//	}
//	//2.把实例返回
//	return &ws
//}
//
////Wallets是对外的，WalletKeyPair是对内的。Wallets调用WalletKeyPair
//func (ws *Wallets)CreateWallet() string {
//	//调用NewWalletKeyPair
//	wallet:=NewWalletKeyPair()
//	//将返回的WalletsKeyPair添加到WalletMap中
//	address:=wallet.GetAddress()
//	ws.WalletsMap[address] = wallet
//
//	res:=ws.SaveToFile()
//	if !res{
//		fmt.Printf("创建钱包失败！\n")
//	}
//	return address
//}
//
//func (ws *Wallets) SaveToFile() bool {
//	var buffer bytes.Buffer
//	gob.Register(elliptic.P256())
//	encoder:=gob.NewEncoder(&buffer)
//
//	err:=encoder.Encode(ws)
//	if err!=nil{
//		fmt.Printf("钱包序列化失败！，err:%v\n",err)
//		return false
//	}
//	content:=buffer.Bytes()
//	err =  ioutil.WriteFile(WalletName,content,0600)
//	if err!=nil{
//		fmt.Printf("钱包创建失败！\n")
//		return false
//	}
//	return true
//}