package main

import (
	"fmt"
	"os"
)

const Usage  = `
      ./bc macc "XXX"  创建商户账户
      ./bc cacc "xxx"  创建顾客账户
      ./bc createbc MACC CACC AMOUT INDEX 创建区块链
      ./bc send MACC CACC AMOUT INDEX " 转账命令"
      ./bc p  打印区块链
      ./bc addbc "*****"  添加数据到区块链
      ./bc ptx 打印交易结果
      ./bc mlist 打印商户钱包地址
      ./bc clist 打印顾客钱包地址
      ./bc get macc 地址 获取地址的余额
      ./bc consensus 共识算法实现
      ./bc test 测试初始化节点
      
`

type CLI struct {

}

func (cli *CLI)Run()  {
	cmds:= os.Args
	if len(cmds)<2{
		fmt.Printf(Usage)
		os.Exit(1)
	}

	switch cmds[1] {
	case "macc":
		fmt.Printf("创建商户账户命令被调用\n")
        cli.CreatMerchant()
	case "cacc":
		fmt.Printf("创建顾客账户命令被调用")
        cli.CreatClient()
	case  "p":
		fmt.Printf("打印区块链命令被调用\n")
		cli.PrintChain()
	case "send":
		fmt.Printf("转账命令被调用\n	")
		if len(cmds)!=6{
			fmt.Printf("send命令发现无效参数，请检查！\n	")
			fmt.Printf(Usage)
			os.Exit(1)
		}
		from:=cmds[2]
		to:=cmds[3]
		amount:=cmds[4]
		index:=cmds[5]
		cli.Send(from,to,amount,index)
	case "ptx":
		fmt.Printf("打印交易结果命令被调用\n")
        cli.PrintTX()
	case "createbc":
		fmt.Printf("创建区块链命令被调用")
		if len(cmds)!=6{
			fmt.Printf("send命令发现无效参数，请检查！\n	")
			fmt.Printf(Usage)
			os.Exit(1)
		}
		from:=cmds[2]
		to:=cmds[3]
		amount:=cmds[4]
		index:=cmds[5]
		cli.CreateBC(from,to,amount,index)
	case "mlist":
		fmt.Printf("打印商户钱包地址命令被调用\n")
		cli.ListMaccs()
	case "clist":
		fmt.Printf("打印顾客钱包地址命令被调用\n")
		cli.ListCaccs()
	case "get":
		fmt.Printf("获取余额命令被调用！\n")
		if len(cmds)!=3{
			fmt.Printf("get命令发现无效参数，请检查！\n	")
			fmt.Printf(Usage)
			os.Exit(1)
		}
		address:=cmds[2]
        cli.GetMacc(address)
	case "consensus":
		fmt.Println("共识算法被调用！\n")
		cli.Consensus()
	case "test":
		fmt.Println("测试初始化节点被调用")
		cli.ConsensNode()
	default:
		fmt.Printf("无效命令，请检查\n ")
		fmt.Printf(Usage)
	}
}