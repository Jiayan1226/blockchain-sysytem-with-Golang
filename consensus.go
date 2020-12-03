package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

/*
1.创建struct
2.初始化struct
3.准备数据
4.共识算法,对于新交易的验证，根据基于区块链应用事先达成的各种验证协议来进行。比如交易的格式，交易的数据结构，格式的语法结构，输入输出、数字签名的正确性等。
5.验证结果
*/

//在共识机制中，
//首先读取文件中的节点
var Consensusnodes []*ConsensusNode

var authenticationSuccess = false
//对外节点


type SafeMap struct {
	Msg map[int64]string
	*sync.RWMutex
}

func NewSafeMap(data map[int64]string) *SafeMap  {
	return &SafeMap{data,&sync.RWMutex{}}
}

//SafeMap 的增删改
func (d *SafeMap)Put(Nodeid int64,command string)  {
	d.Lock()
	defer d.Unlock()
	d.Msg[Nodeid]=command
}

func (d *SafeMap)Length() int  {
	d.RLock()
	defer d.RUnlock()
	value:=len(d.Msg)
	return value
}



type ConsensusNode struct {
	//每个节点虚拟地址
	Nodeid int64
	//每个节点的消息池(节点名称，以及是否收到信息flag)
	MsgMap *SafeMap
	//自己收到的交易信息
   Tx Transaction
	*sync.RWMutex
}

func  NewConsensusNode(id int64) *ConsensusNode {
	tx:=Transaction{}
	//首先需要给map分配空间
	msg:= make(map[int64]string)
	return &ConsensusNode{id,NewSafeMap(msg),tx,&sync.RWMutex{}}
}
//初始化共识结点
func LoadMerchant()  {
	//1.加载文件
	if !IsFileExist(MerFile){
		fmt.Printf("钱包文件不存在，准备创建!\n")
	}
	content,err:= ioutil.ReadFile(MerFile)
	if err!=nil{
		log.Panic(err)
	}
	gob.Register(elliptic.P256())
	decoder:=gob.NewDecoder(bytes.NewReader(content))
	var mas Maccs

	err = decoder.Decode(&mas)
	if err!=nil{
		log.Panic(err)
	}
	//2.取出文件中的macc用户id号
	for id,_:=range mas.Maccsmap{
		conode:= NewConsensusNode(id)
		Consensusnodes = append(Consensusnodes,conode)
	}
}

func MatchNode(from string) *ConsensusNode{
	LoadMerchant()
	mon:=ConsensusNode{}
	for i,con:=range Consensusnodes{
      if con.Nodeid == StoInt64(from){
      	mon = *Consensusnodes[i]
	   }
	}
   return &mon
}
//程序开始之前首先读取文件中存在的节点数量
//提供加密对称加密之后的密文，需要验证签名，不需要回溯前区块（可节省时间）
func (con *ConsensusNode)broadcast(a string)  {
     //遍历所有节点进行广播
     con.Lock()
     defer con.Unlock()
     for _,conode:=range Consensusnodes{
     	if conode.Nodeid == con.Nodeid{
     		continue
		}
     	//广播的内容是自己节点编号，以及自己是否收到信息的flag
		  //conode.MsgMap.Msg[con.Nodeid] = a
		  conode.MsgMap.Put(con.Nodeid,a)
        fmt.Println(conode.Nodeid," :  ",conode.MsgMap)
	  }
}

//不需要算力的计算，但是需要
func (con *ConsensusNode)authentication(txs []*Transaction,length int,bc *Blockchain)  {
	con.Lock()
	defer con.Unlock()
	fmt.Println("收到广播信息，现在开始验证：")
	value:=con.MsgMap.Length()
	if value > length/3{
		//2,验证签名是否有效
		//验证签名，以及上一个区块的prehash
		var hashes []byte
		for _,tx:=range txs{
			if !con.VerifyTransaction(tx){
				fmt.Println("发现无效交易，验证终止！")
				break
			}
			txid/*[]byte*/:=tx.Txid
			hashes=append(hashes,txid...)
		}
		hash:=append(hashes[:])
		if !bytes.Equal(hash,bc.tail){
           fmt.Println("发现无效区块！")
		}
		fmt.Println("区块验证成功！")
	}
}


func (con *ConsensusNode)SignTransaction(tx *Transaction,privateKey *ecdsa.PrivateKey)  {
	con.Lock()
	defer con.Unlock()
	tx.Sign(privateKey)
}

func (con *ConsensusNode)VerifyTransaction(tx *Transaction) bool {
	con.Lock()
	con.Unlock()
	//所有交易必须全部验证
	return tx.Verify()
}
func (con *ConsensusNode)commit()  {
  //验证成功后，每个节点对验证成功的信息进行汇总
}


func ConsensusMechanism(txs []*Transaction,bc *Blockchain,xf1 *sync.WaitGroup) {
	defer xf1.Done()
	LoadMerchant()
	length := len(Consensusnodes)
	xf := &sync.WaitGroup{}
	xf.Add(len(Consensusnodes) * 2)
	for i, conode := range Consensusnodes {

		//广播收到消息，确认后验证、验证完再广播
		// defer xf.Done()
		fmt.Println("第",i,"次")
		go conode.broadcast("prepare")
		go conode.authentication(txs, length, bc)
 		//打包结点对收到消息进行汇总，3/4打包
		xf.Wait()
		time.Sleep(time.Second * 2)
		fmt.Println("共识校验结束")
	}

}


var globalQuit chan struct{}

var tss *threadSafeSlice

//当监听者数量未知时：
type threadSafeSlice struct {
	sync.Mutex
	workers []*worker
}

func (slice *threadSafeSlice)Push(w *worker)  {
	slice.Lock()
	defer slice.Unlock()

	slice.workers = append(slice.workers,w)
}
func (slice *threadSafeSlice)Iter(routine func(w *worker))  {
	slice.Lock()
	defer slice.Unlock()

	for _,worker:=range slice.workers{
		routine(worker)
	}
}

//创建每一个监听节点
type worker struct {
	name   string //ID码
	Broadcount int64
	Source chan interface{}//专有的广播通道
	Quit   chan struct{}// 关闭广播通道
	AuthenticQueuew  chan []*Transaction  //交易的共识通道，即来即走
	Finished    chan  bool//共识之后的通知通道，无缓冲
    ConsensusQueue chan chan []*Transaction                  //等待分配交易的节点的
}

func (w *worker)Start()  {
	w.Source = make(chan interface{})
	go func() {
		for {
			select {
			case msg:=<-w.Source://广播内容
			     fmt.Println("=========>",w.name,"收到信息，信息是：",msg,"接下来将进行广播！","目前收到信息次数：",w.Broadcount)
			    // w.BroadCast()
			case <- w.Quit://关闭通道
			     fmt.Println(w.name,"quit!")
				 return
			}
		}
	}()
}

func (w *worker)BroadCast()  {
	count:= 0
	msg := "broadcast"
	var broadcastMsg string
	go func() {
		for {
			select{
			case <- globalQuit:
				fmt.Println("请求退出程序信息！")
				return
			case <- time.Tick(1 * time.Second):
				count++
				broadcastMsg = fmt.Sprintf("%s-%d",msg,count)
                fmt.Println("Broadcast message is",broadcastMsg)
				tss.Iter(func(w *worker) {
					w.Source <- SendMsg
				})
			}
		}
	}()
}
//打包节点派发消息
//TODO
func SendMsg()  {
	globalQuit = make(chan struct{})
	tss = &threadSafeSlice{}
	go func() {
		msg:="true"
		count:=0
		var sendMsg string
		for {
			select {
			case <-globalQuit:
				fmt.Println("Stop send message")
				return
			case <-time.Tick(5 *time.Second)://要充分考虑系统运行的时间
				count++
				sendMsg = fmt.Sprintf("%s-%d",msg,count)//输入信息
				fmt.Println("Send message is",sendMsg)
				tss.Iter(func(w *worker) {
					//w.Broadcount = 0
					w.Source <- sendMsg //进行广播给每一个
					w.Broadcount++ //收到信息次数进行+1
				})
			}
		}
	}()
}
