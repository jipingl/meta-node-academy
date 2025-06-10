package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jipingl/demo/erc20"
	"github.com/jipingl/demo/util"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	EndPoint string `envconfig:"END_POINT"`
	FromAddrPubKey string `envconfig:"FROM_ADDRESS_PUBLIC_KEY"`
	FromAddrPrvKey string `envconfig:"FROM_ADDRESS_PRIVATE_KEY"`
	ToAddrPubKey  string `envconfig:"TO_ADDRESS_PUBLIC_KEY"`
	ToAddrPrvKey string `envconfig:"TO_ADDRESS_PRIVATE_KEY"`
}
var config Config

func task_1() {
	// 连接测试网
	client, err := ethclient.Dial(config.EndPoint)
	util.Err(err)
	// 查询区块信息
	block, err := client.BlockByNumber(context.Background(), big.NewInt(1))
	util.Err(err)
	// 打印区块信息
	fmt.Println("Number: ", block.Number())
	fmt.Println("Hash: ", block.Hash().Hex())
	fmt.Println("Timestamp: ", block.Time())
	fmt.Println("TxCount: ", block.Transactions().Len())
	fmt.Println("GasUsed: ", block.GasUsed())
	
	fromAddress := common.HexToAddress(config.FromAddrPubKey)
	toAddress := common.HexToAddress(config.ToAddrPubKey)
	amount := big.NewInt(1000000000000000000) // 1ETH

	// Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	util.Err(err)

	// GasLimit
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To: &toAddress,
		Value: amount,
	})
	util.Err(err)
	fmt.Println(gasLimit)

	// GasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	util.Err(err)

	// ChainId
	chainId, err := client.ChainID(context.Background())
	util.Err(err)

	// 转账交易
	fromPrvECDSA, err := crypto.HexToECDSA(config.FromAddrPrvKey)
	util.Err(err)
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, []byte{})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), fromPrvECDSA)
	util.Err(err)
	err = client.SendTransaction(context.Background(), signedTx)
	util.Err(err)

	// 打印交易哈希
	fmt.Println("tx sent: ", signedTx.Hash().Hex())
}

func task_2() {
	// 连接测试网
	client, err := ethclient.Dial(config.EndPoint)
	util.Err(err)

	// 调用者
	fromAddress := common.HexToAddress(config.FromAddrPubKey)

	// Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	util.Err(err)

	// GasLimit
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To: nil,
		Data: common.FromHex(erc20.MyTokenBin),
	})
	util.Err(err)

	// GasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	util.Err(err)

	// ChainId
	chainId, err := client.ChainID(context.Background())
	util.Err(err)

	// 签名器
	prvKeyECDSA, err := crypto.HexToECDSA(config.FromAddrPrvKey)
	util.Err(err)
	auth := bind.NewKeyedTransactor(prvKeyECDSA, chainId)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = common.Big0
	auth.GasPrice = gasPrice
	auth.GasLimit = gasLimit
	// 部署
	address, tx, contract, err := erc20.DeployMyToken(auth, client)
	util.Err(err)
	fmt.Println("address: ", address)
	fmt.Println("tx: ", tx.Hash().Hex())



	// 铸造代币
	auth0 := bind.NewKeyedTransactor(prvKeyECDSA, chainId)
	auth0.Value = big.NewInt(100)
	auth0.GasPrice = gasPrice
	// nonce
	nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
	util.Err(err)
	auth0.Nonce = big.NewInt(int64(nonce))
	// 执行合约函数
	tx, err = contract.Mint(auth0)
	util.Err(err)
	fmt.Println("mint tx hash: ", tx.Hash().Hex())
	// 等待被挖矿
	receipt, err := bind.WaitMined(context.Background(), client, tx.Hash())
	util.Err(err)
	if receipt.Status == 0 {
		panic("tx failed")
	}


	// 查询余额
	balance, err := contract.BalanceOf(&bind.CallOpts{
		Pending: true,
		From: fromAddress,
	}, fromAddress)
	util.Err(err)
	fmt.Println("balance shouled be 10000: ", balance)
}

func main()  {
	// 加载.env配置文件
	err := godotenv.Load()
	util.Err(err)
	// 绑定环境变量到结构体
	err = envconfig.Process("", &config)
	util.Err(err)
	// 连接测试网查询区块和交易信息
	// task_1()
	// 合约代码的生成与调用
	task_2()
}