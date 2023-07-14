package main

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/SideSwapIN/Analystic/internal/db"
	"github.com/SideSwapIN/Analystic/internal/initialize"
	"github.com/SideSwapIN/Analystic/internal/logger"
	"github.com/SideSwapIN/Analystic/internal/model"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ChainInfo struct {
	StartBlock        *big.Int
	ChainID           int64
	RPC               string
	BlockTimeInterval int64
	Deviation         int64
	Router            common.Address
}

var (
	QueryTimeInterval = time.Second * 1
	cacheOldBlock     = "CACHES:LISTEN:OLD_BLOCK"
	OneBigNumber      = big.NewInt(1)
	chains            = []ChainInfo{
		// {
		// 	ChainID:           56,
		// 	RPC:               "https://rpc.ankr.com/bsc",
		// 	BlockTimeInterval: 3000,
		// 	Deviation:         1000,
		// 	Router:            common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E"),
		// },
		{
			StartBlock:        big.NewInt(166047),
			ChainID:           51178,
			RPC:               "https://pre-alpha-us-http-geth.opside.network",
			BlockTimeInterval: 200,
			Deviation:         0,
			Router:            common.HexToAddress("0xAAb8FCD8DD22a5de73550F8e67fF9Ca970d1257E"),
		},
		{
			StartBlock:        big.NewInt(29807),
			ChainID:           12008,
			RPC:               "https://pre-alpha-zkrollup-rpc.opside.network/public",
			BlockTimeInterval: 200,
			Deviation:         0,
			Router:            common.HexToAddress("0x7A9a466DE08747bC8Ad79aBA6D8dCE9D64E5C767"),
		},
		// {
		// 	ChainID:           56,
		// 	RPC:               "https://rpc.ankr.com/bsc",
		// 	BlockTimeInterval: 3000,
		// 	Deviation:         1000,
		// 	Router:            common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E"),
		// },
		// 51178: "https://pre-alpha-us-http-geth.opside.network",
		// 12008: "https://pre-alpha-zkrollup-rpc.opside.network/public",
	}

	MethodMap = map[string]model.RouterMethod{
		"e8e33700": model.AddLiquidityRouterMethod, // addLiquidity
		"f305d719": model.AddLiquidityRouterMethod, // addLiquidityETH

		"baa2abde": model.RemoveLiquidityRouterMethod, // removeLiquidity
		"02751cec": model.RemoveLiquidityRouterMethod, // removeLiquidityETH
		"af2979eb": model.RemoveLiquidityRouterMethod, // removeLiquidityETHSupportingFeeOnTransferTokens
		"ded9382a": model.RemoveLiquidityRouterMethod, // removeLiquidityETHWithPermit
		"5b0d5984": model.RemoveLiquidityRouterMethod, // removeLiquidityETHWithPermitSupportingFeeOnTransferTokens
		"2195995c": model.RemoveLiquidityRouterMethod, // removeLiquidityWithPermit

		"fb3bdb41": model.SwapRouterMethod, // swapETHForExactTokens
		"7ff36ab5": model.SwapRouterMethod, // swapExactETHForTokens
		"b6f9de95": model.SwapRouterMethod, // swapExactETHForTokensSupportingFeeOnTransferTokens
		"18cbafe5": model.SwapRouterMethod, // swapExactTokensForETH
		"791ac947": model.SwapRouterMethod, // swapExactTokensForETHSupportingFeeOnTransferTokens
		"38ed1739": model.SwapRouterMethod, // swapExactTokensForTokens
		"5c11d795": model.SwapRouterMethod, // swapExactTokensForTokensSupportingFeeOnTransferTokens
		"4a25d94a": model.SwapRouterMethod, // swapTokensForExactETH
		"8803dbee": model.SwapRouterMethod, // swapTokensForExactTokens
	}
)

func main() {
	initialize.Init()
	for _, chainInfo := range chains {
		go Start(chainInfo)
	}
	select {}
	// if block == nil {
	// 	latest, err := client.BlockByNumber(context.Background(), nil)
	// 	if err != nil {
	// 		panic(fmt.Errorf("client.BlockByNumber RPC: %s, error: %v", config.GetConfig().EVM.RPC, err))
	// 	}
	// 	block = latest
	// }
}

func Start(chainInfo ChainInfo) {
	cacheKey := fmt.Sprintf("%s:%d", cacheOldBlock, chainInfo.ChainID)
	client, err := ethclient.Dial(chainInfo.RPC)
	if err != nil {
		panic(fmt.Errorf("ethclient.Dial RPC: %s, error: %v", chainInfo.RPC, err))
	}
	var oldBlockBigInt *big.Int
	oldBlockStr, err := db.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		oldBlockInt, err := strconv.ParseInt(oldBlockStr, 10, 64)
		if err == nil {
			oldBlockBigInt = big.NewInt(oldBlockInt)
		}
	}
	if (chainInfo.StartBlock != nil && oldBlockBigInt == nil) ||
		(chainInfo.StartBlock != nil && chainInfo.StartBlock.Cmp(oldBlockBigInt) > 0) {
		oldBlockBigInt = chainInfo.StartBlock
	}
	// query := ethereum.FilterQuery{
	// 	Addresses: []common.Address{},
	// }
	router := chainInfo.Router.Hex()
	for {
		time.Sleep(time.Microsecond * 200)
		optionID := time.Now().UnixMilli()
		var blockBigInt *big.Int
		if oldBlockBigInt != nil {
			blockBigInt = big.NewInt(0).Add(oldBlockBigInt, OneBigNumber)
		}

		block, err := client.BlockByNumber(context.Background(), blockBigInt)
		if err != nil {
			logger.Debugf("ethclient.Dial RPC: %s, block: %s, error: %v\n", chainInfo.RPC, blockBigInt.String(), err)
			continue
		}

		txs := block.Transactions()
		senderOperations := make([]model.SenderOperation, 0)
		var sw sync.WaitGroup
		for _, tx := range txs {
			sw.Add(1)
			go func(tx *types.Transaction) {
				defer sw.Done()
				if tx.To() != nil && tx.To().Hex() == router && len(tx.Data()) >= 4 {
					methodSig := common.Bytes2Hex(tx.Data()[:4])
					if methodType, ok := MethodMap[methodSig]; ok {
						from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
						if err != nil {
							logger.Dump(tx)
							logger.Errorf("types.Sender error: %v", err)
							return
						}
						senderOperations = append(senderOperations, model.SenderOperation{
							To:          router,
							From:        from.Hex(),
							TxHash:      tx.Hash().Hex(),
							ChainID:     chainInfo.ChainID,
							BlockNumber: block.NumberU64(),
							BlockTime:   block.Time(),
							Type:        methodType,
						})
					}
				}
			}(tx)
		}
		sw.Wait()
		if len(senderOperations) > 0 {
			err = model.CreateSenderOperations(senderOperations)
			if err != nil {
				logger.Errorf("model.CreateSenderOperations error: %v", err)
				continue
			}
		}

		if chainInfo.BlockTimeInterval > chainInfo.Deviation {
			devia := chainInfo.BlockTimeInterval - chainInfo.Deviation - (time.Now().UnixMilli() - int64(block.Time()*1000))
			if devia > 0 {
				time.Sleep(time.Duration(devia) * time.Millisecond)
			}
		}
		oldBlockBigInt = block.Number()
		db.RedisClient.Set(context.Background(), cacheKey, oldBlockBigInt.String(), 0)
		logger.Debugf("[%d]listenChain block: %d chainID: %d end", optionID, oldBlockBigInt, chainInfo.ChainID)
	}
}

func ScanBlocks(client *ethclient.Client, chianInfo ChainInfo, query ethereum.FilterQuery, startBlock, endBlock *big.Int, optionID int64) {

}
