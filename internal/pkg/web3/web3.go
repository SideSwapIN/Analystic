package web3

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/SideSwapIN/Analystic/internal/logger"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

func SendTransaction(
	client *ethclient.Client,
	from common.Address,
	recipient common.Address,
	contract common.Address,
	private *ecdsa.PrivateKey) (*types.Transaction, error) {
	// client, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545") // BSC TESTNET

	data, err := packTranserAllTokenData(client, from, recipient, contract)
	if err != nil {
		logger.Errorf("service.UserService sendTransaction packTranserAllTokenData error: %v", err)
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		logger.Errorf("service.UserService sendTransaction PendingNonceAt error: %v", err)
		return nil, err
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		logger.Errorf("service.UserService sendTransaction ChainID error: %v", err)
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Errorf("service.UserService SuggestGasPrice ChainID error: %v", err)
		return nil, err
	}

	msg := ethereum.CallMsg{
		From: from,
		To:   &contract,
		Data: data,
	}
	estimateGas, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		logger.Errorf("service.UserService sendTransaction EstimateGas error: %v", err)
		return nil, err
	}
	signer := types.LatestSignerForChainID(chainID)
	LegacyTx := &types.LegacyTx{
		Nonce:    nonce,
		To:       &contract,
		Value:    big.NewInt(0),
		Data:     data,
		Gas:      estimateGas,
		GasPrice: gasPrice,
	}
	tx, err := types.SignNewTx(private, signer, LegacyTx)
	if err != nil {
		logger.Errorf("service.UserService SignNewTx ChainID error: %v", err)
		return nil, err
	}
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		logger.Errorf("service.UserService SendTransaction SendTransaction error: %v", err)
		gasFee := decimal.NewFromBigInt(gasPrice, -18).Mul(decimal.NewFromInt(int64(estimateGas)))
		return nil, fmt.Errorf("%s, gas fee >= %s BNB", err.Error(), gasFee.String())
	}
	return tx, nil
}

func packTranserAllTokenData(client *ethclient.Client, from common.Address, recipient common.Address, contract common.Address) ([]byte, error) {

	amount, err := GetTokenBalance(client, contract, from)
	if err != nil {
		return nil, err
	}
	logger.Debugf("packTranserAllTokenData amount: %v", amount)
	balanceOfFunc := crypto.Keccak256Hash([]byte("transfer(address,uint256)")).Hex()[:10]
	d, err := hexutil.Decode(balanceOfFunc)
	if err != nil {
		logger.Errorf("GetTokenBalance hexutil.Decode err: %v", err)
		return nil, err
	}

	agrs := abi.Arguments{
		abi.Argument{
			Type: Address,
		},
		abi.Argument{
			Type: Uint256,
		},
	}
	logger.Dump("balanceOf: ", amount)
	data, err := agrs.Pack(recipient, amount)
	if err != nil {
		logger.Errorf("GetTokenBalance agrs.Pack err: %v", err)
		return nil, err
	}
	d = append(d, data...)
	return d, nil
}

var (
	Uint256, _ = abi.NewType("uint256", "", nil)
	Address, _ = abi.NewType("address", "", nil)
)

func GetTokenBalance(client *ethclient.Client, token common.Address, walletAddress common.Address) (*big.Int, error) {
	balanceOfFunc := crypto.Keccak256Hash([]byte("balanceOf(address)")).Hex()[:10]
	d, err := hexutil.Decode(balanceOfFunc)
	if err != nil {
		logger.Errorf("GetTokenBalance hexutil.Decode err: %v", err)
		return nil, err
	}

	agrs := abi.Arguments{
		abi.Argument{
			Type: Address,
		},
	}
	data, err := agrs.Pack(walletAddress)
	if err != nil {
		logger.Errorf("GetTokenBalance agrs.Pack err: %v", err)
		return nil, err
	}
	d = append(d, data...)
	msg := ethereum.CallMsg{
		To:   &token,
		Data: d,
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}
	logger.Dump(result)
	balance := common.BytesToHash(result).Big()
	return balance, nil
}
