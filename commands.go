package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func (cli *CLI)PrintChain()  {
	bc:=NewBlockchain()
	if bc == nil {
		return
	}
	defer bc.db.Close()
	it:= bc.NewIterator()
	for{
		block:=it.Next()
		fmt.Printf("++++++++++++++++++++++\n")
		fmt.Printf("Version : %d\n",block.Version)
		fmt.Printf("Hash : %x\n",block.Hash)
		fmt.Printf("PreBlockHash : %x\n",block.PrevBlockHash)
		fmt.Printf("MerkleRoot : %x\n",block.MerKleRoot)
		timeFormat:=time.Unix(int64(block.Timestamp),0).Format("2006-01-02 15:04:05")
		fmt.Printf("Timestamp : %s\n",timeFormat)
		fmt.Printf("Transaction : %x\n",block.Transactions[0].Txid)
		fmt.Printf("\n")
		//TODO
		//共识算法调用
		if bytes.Equal(block.PrevBlockHash,[]byte{}){
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}
}

func (cli *CLI)PrintTX()  {
    bc:=NewBlockchain()
    if bc == nil{
		return
	}
    defer bc.db.Close()
    it := bc.NewIterator()

    for{
    	block:= it.Next()

    	fmt.Printf("\n+++++++ 新的区块 +++++++++\n")
    	for _,tx:= range block.Transactions{
    		fmt.Printf("tx:%v\n",tx)//打印一个结构
		}
    	if len(block.PrevBlockHash)==0{
    		break
		}
	}
}

func (cli *CLI)CreatMerchant()  {
	maccs:=NewMaccs()
	address:=maccs.Create()
	fmt.Printf("新的商户地址为：%d",address)
}
func (cli *CLI)CreatClient()  {
	caccs:=NewCaccs()
	address:=caccs.Create()

	fmt.Printf("新的顾客地址为：%d",address)
}

//把商户的地址持久化
func (cli *CLI)ListMaccs()  {
	maccs:=NewMaccs()
	addresses:= maccs.ListMaccs()
	for _,address:= range addresses{
		fmt.Printf("address:  %d\n",address)
	}
}
func (cli *CLI)ListCaccs()  {
	caccs:=NewCaccs()
	addresses:= caccs.ListCaccs()
	for _,address:= range addresses{
		fmt.Printf("address:  %d\n",address)
	}
}
func (cli *CLI)CreateBC(from string,to string,amount string,index string)  {
	maccs:= NewMaccs()
	macc:= maccs.Maccsmap[StoInt64(from)]
	//macc.GetMacc()
	caccs:=NewCaccs()
	cacc:=caccs.Caccsmap[StoInt64(to)]
	bc:=CreateBlockchain(cacc,macc,StoFloat64(amount),uint64(StoInt64(index)))
	if bc!=nil{
		defer bc.db.Close()
	}
	fmt.Printf("区块链创建成功！")
}

func (cli *CLI)Send(from string,to string,amount string,index string)  {
	runtime.GOMAXPROCS(4)
	t:=time.Now()
	//初始化节点
	xf1:=&sync.WaitGroup{}
	xf1.Add(1)
    //主节点实例化
    //开始进行交易
	//cacc.GetCacc()
	bc:=NewBlockchain()
	if bc ==nil{
		return
	}
	defer bc.db.Close()
	txs:=[]*Transaction{}
	for i:=1;i<500;i++{
		fmt.Println("+++++++++++++++++第",i,"次交易：+++++++++++++++++")
		tx:=NewTransaction(from,to,StoFloat64(amount),uint64(StoInt64(index)),bc)
		//其他节点开始共识，
		if tx!=nil{
			txs = append(txs,tx)
		}else{
			fmt.Printf("发现无效交易，过滤！\n")
		}
	}
	bc.AddBlock(MatchNode(from),txs)
	//产生交易后，节点开始广播。
	/*进行共识*/
    go ConsensusMechanism(txs,bc,xf1)
   //进行共识
    //xf1.Wait()
	//
	//fmt.Printf("交易后：")
	elapsed:=time.Since(t)
	fmt.Println("app elapsed:",elapsed)
	//macc.GetMacc()
	//cacc.GetCacc()
}

func (cli *CLI)GetMacc(address string)  {
	maccs:=NewMaccs()
	macc:=maccs.Maccsmap[StoInt64(address)]
	macc.GetMacc()
}

func (cli *CLI)Consensus(){
	//txs:=*[]Transaction{}


	time.Sleep(time.Second * 120)
	/*1.主线程模拟用户完成交易，
	           放入交易池，
	           并向所有的节点广播
	           同时交易池中设定一个时间，时间完成后将交易池复制，
	           交易池清空，同时打包结点对上一段时间完成的交易进行打包验证。
	  2.每个节点开启一个线程，
	     线程的工作：
	     2.1 首先判断是否接收到信息，
	           如果接收到广播
	               回复ok，
	          没有接收到或者作恶节点广播
	               不广播或者回复false
	    2.2 判断接受信息的数量，如果超过总数量的1/3，进行校验，校验成功广播校验结果
	    线程中加入加入时间延迟。如果超出时间未广播，将不考虑该节点
	  主节点将交易进行打包，放入数据库。
	*/

}
//TODO
//激励结算模型，
func Incentive()  {
	//通过结算优惠券的过程，进行激励机制
	//在论文中描述，之后再做实验
}

func (cli *CLI)Test()  {
     fmt.Println("请问是否需要继续交易：")
     for{
     	/*
     	1.持续监听客户端是否传来交易
     	2.
     	*/
     	time.Sleep(time.Second * 60)
	 }

}
func (cli *CLI)ConsensNode()  {
	//1秒钟添加一个新的worker至slice中

	go func() {
		name:="worker"
		for i:=0;i<5;i++{
			time.Sleep(1*time.Second)
			w:=&worker{
				Broadcount:0,
				name:fmt.Sprintf("%s%d",name,i),
				Quit:globalQuit,
			}
			w.Start()
			tss.Push(w)
		}
	}()
	go SendMsg()

	//截获退出信号
	c:=make(chan os.Signal,1)
	signal.Notify(c,syscall.SIGINT,syscall.SIGTERM)
	for sig:=range c{
		switch sig {
		case syscall.SIGINT,syscall.SIGTERM://获取退出信号时，关闭globalQuit，让所有的监听者退出
		     close(globalQuit)
		     time.Sleep(1*time.Second)
		     return
		}
	}
}
