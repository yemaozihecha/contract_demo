package contract_demo

import (
	token "contract_demo/solidity"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type TokenManager struct {
	tokenAddr common.Address
	client    *ethclient.Client
	instance  *token.Solidity
}

// NewERC20TokenManager 实例一个solidity的token管理
func NewERC20TokenManager(contractAddr string, client *ethclient.Client) (tokenManager *TokenManager, err error) {
	tokenAddr := common.HexToAddress(contractAddr)
	instance, err := token.NewSolidity(tokenAddr, client)
	if err != nil {
		return
	}
	return &TokenManager{
		tokenAddr: tokenAddr,
		client:    client,
		instance:  instance,
	}, nil
}

func (tm *TokenManager) BalanceAt(accountAddr string) (bal *big.Int, err error){
	addr := common.HexToAddress(accountAddr)
	return tm.instance.BalanceOf(&bind.CallOpts{}, addr)
}

// Name 获取代币名称
func (tm *TokenManager) Name() (string, error){
	return tm.instance.Name(&bind.CallOpts{})
}

// Symbol 获取代币符号
func (tm *TokenManager) Symbol() (string, error){
	return tm.instance.Symbol(&bind.CallOpts{})
}

// Decimals 获取代币小数点后有多少位
func (tm *TokenManager) Decimals() (uint8, error){
	return tm.instance.Decimals(&bind.CallOpts{})
}