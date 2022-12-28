package contract_demo

import (
	"context"
	"contract_demo/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
	"time"
)

var (
	chainAddr = "http://127.0.0.1:7545"
	privateKeyStr = "ec8a2bbd238fcafdcac3d3b01c288506001e2344da915b25ccd1ff8b1972218d"
	accountA = common.HexToAddress("0x57Ec540D293f2bd8c60Fa50A31Faee141F10a1A0") // 部署者
	accountB = common.HexToAddress("0x437Da62699b8A161f468d74cbb66f1dF5EB48776") //

	// 合约地址
	contractAddr = common.HexToAddress("0x5A5Dac2346f6B885b067c242bF017A1aDC1b014e")
)

func TestERC20Manager_BalanceAt(t *testing.T) {
	m := NewERC20Manager(chainAddr)
	if err := m.Connect(); err != nil {
		panic(err)
	}
	defer m.Close()
	balance, err := m.BalanceAt(context.Background(), accountA, nil)
	if err != nil {
		fmt.Printf("balanceAt err %s \n", err.Error())
		return
	}
	fmt.Printf("account %s balance is %v \n", accountA.String(), utils.ToETH(balance))
}

func TestTransaction(t *testing.T) {
	m := NewERC20Manager(chainAddr)
	if err := m.Connect(); err != nil {
		panic(err)
	}
	defer m.Close()
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	privateKey, fromAddr, err := utils.GetAddrFromPrivateKey(privateKeyStr)
	if err != nil {
		fmt.Printf("get addr from private key err is %s", err.Error())
		return
	}
	fmt.Println("fromAddr", fromAddr)

	ercToken, err := NewERC20TokenManager(contractAddr.String(), m.client)
	if err != nil {
		fmt.Printf("NewERC20TokenManager err %v", err)
		return
	}
	balanceToken, err := ercToken.BalanceAt(accountA.String())
	if err != nil {
		fmt.Printf("BalanceAt err %v", err)
		return
	}
	fmt.Println("accountA balanceToken", balanceToken)

	nonce, err := m.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		fmt.Printf("PendingNonceAt err is %s", err.Error())
		return
	}
	// 代币传输不需要ETH
	value := big.NewInt(0)
	gasPrice, err := m.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Printf("SuggestGasPrice err is %s", err.Error())
		return
	}
	methodID := utils.GetMethodID([]byte("transfer(address,uint256)"))
	// 转移目标地址
	// 需要向左填充32位
	paddedAddress := common.LeftPadBytes(accountB.Bytes(), 32)

	// 转移数量， 10个代币
	amount := new(big.Int)
	amount.SetString("10", 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// 使用EstimateGas 估算的值一直太小导致，out of gas
	//gasLimit, err := m.EstimateGas(ctx, accountB, data)
	//if err != nil {
	//	fmt.Printf("EstimateGas err is %s", err.Error())
	//	return
	//}

	gasLimit := uint64(6721975)
	//构建交易事务
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &contractAddr,
		Value:    value,
		Data:     data,
	})
	chainID, err := m.NetworkID(ctx)
	if err != nil {
		fmt.Printf("get NetworkID err is %s", err.Error())
		return
	}
	// 需要发件人的私钥对事务进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Printf("get SignTx err is %s", err.Error())
		return
	}
	// 广播交易
	err = m.SendTransaction(ctx, signedTx)
	if err != nil {
		fmt.Printf("SendTransaction err is %s", err.Error())
		return
	}
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
