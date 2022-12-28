package contract_demo

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type ERC20Manager struct {
	client    *ethclient.Client
	chainAddr string
}

func NewERC20Manager(chainAddr string) *ERC20Manager {
	return &ERC20Manager{
		chainAddr: chainAddr,
	}
}

func (m *ERC20Manager) Connect() (err error){
	if m.chainAddr == "" {
		return fmt.Errorf("chain addr is emtpy")
	}
	m.client, err = ethclient.Dial(m.chainAddr)
	return
}

func (m *ERC20Manager) Close() {
	m.client.Close()
}

// BalanceAt 查询某个账户地址的余额
// blockNumber 传区块号能让您读取该区块时的账户余额
func (m *ERC20Manager) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	return m.client.BalanceAt(ctx, account, blockNumber)
}

// PendingBalanceAt 待处理的账户余额是多少
func (m *ERC20Manager) PendingBalanceAt(ctx context.Context, account common.Address) (balanceStr *big.Int, err error) {
	return m.client.PendingBalanceAt(context.Background(), account)
}

// IsContractAddr 判断该地址是否是合约地址
func (m *ERC20Manager)IsContractAddr(address common.Address) (isContract bool, err error) {
	byteCode, err := m.client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return
	}
	return len(byteCode) > 0, nil
}

// PendingNonceAt 返回应该使用的下一个nonce
func (m *ERC20Manager) PendingNonceAt(ctx context.Context, fromAddress common.Address)(uint64, error) {
	return m.client.PendingNonceAt(ctx, fromAddress)
}

//SuggestGasPrice 据'x'个先前块来获得平均燃气价格
func (m *ERC20Manager) SuggestGasPrice (ctx context.Context) (*big.Int, error){
	return m.client.SuggestGasPrice(ctx)
}

// EstimateGas 估算所需要的gas值
func (m *ERC20Manager) EstimateGas(ctx context.Context, toAddr common.Address, data []byte )(uint64, error) {
	return m.client.EstimateGas(ctx, ethereum.CallMsg{
		//From:       common.HexToAddress("0x57Ec540D293f2bd8c60Fa50A31Faee141F10a1A0"),
		To:         &toAddr,
		Data:       data,
	})
}

// NetworkID 获取链id
func (m *ERC20Manager) NetworkID(ctx context.Context) (*big.Int, error){
	return m.client.NetworkID(ctx)
}

func (m *ERC20Manager) SendTransaction (ctx context.Context, tx *types.Transaction) error {
	return m.client.SendTransaction(ctx, tx)
}
