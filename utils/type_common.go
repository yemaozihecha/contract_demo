package utils

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"math"
	"math/big"
)

// ToETH wei转换eth值
func ToETH(wei *big.Int) *big.Float {
	fBalance := new(big.Float)
	fBalance.SetString(wei.String())
	return new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
}

// GetAddrFromPrivateKey 通过私钥派生出地址
func GetAddrFromPrivateKey(hexKey string) (privateKey *ecdsa.PrivateKey, address common.Address, err error){
	privateKey, err = crypto.HexToECDSA(hexKey)
	if err != nil {
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA)
	return
}

// GetMethodID 获取方法id
func GetMethodID(fnSign []byte)(methodID []byte) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(fnSign)
	methodID = hash.Sum(nil)[:4]
	return
}