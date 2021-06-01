package mdexswap

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wo4zhuzi/zujian/mdexswap/contracts"
	"github.com/wo4zhuzi/zujian/utils"
)

type SwapRouter struct {
	Client       *ethclient.Client
	SwapContract *contracts.SwapRouter
	Address      string
	Factory      string
	Farm         *MdexFarm
}

func NewSwapRouter(address string, client *ethclient.Client, Farm *MdexFarm) (*SwapRouter, error) {
	swapRouterContract, err := contracts.NewSwapRouter(common.HexToAddress(address), client)
	if err != nil {
		return nil, err
	}
	factory, _ := swapRouterContract.Factory(&bind.CallOpts{})
	return &SwapRouter{
		Factory:      factory.String(),
		Address:      address,
		Client:       client,
		SwapContract: swapRouterContract,
		Farm:         Farm,
	}, nil
}

func (c *SwapRouter) SwapExactTokenTo(fromToken, toToken string, sendAmount, amountMin *big.Int) (*types.Transaction, error) {
	auth, err := c.Farm.TokenBasic.CreateTx()
	if err != nil {
		return nil, err
	}
	wallet := c.Farm.FarmConfig.Wallet
	amountIn := sendAmount
	gasPrice, err := c.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	ABI, err := abi.JSON(strings.NewReader(contracts.SwapRouterABI))
	if err != nil {
		return nil, err
	}

	nonce, err := c.Client.PendingNonceAt(context.Background(), common.HexToAddress(wallet))
	if err != nil {
		return nil, err
	}
	t := time.Now()
	deadLine := t.Add(time.Hour * 24).Unix()
	toContract := common.HexToAddress(c.Address)
	txData, err := ABI.Pack("swapExactTokensForTokens", amountIn, amountMin, []common.Address{
		common.HexToAddress(fromToken),
		common.HexToAddress(toToken),
	}, common.HexToAddress(wallet), big.NewInt(deadLine))
	if err != nil {
		return nil, err
	}
	gas, err := c.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     common.HexToAddress(wallet),
		To:       &toContract,
		GasPrice: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		Value:    utils.ToWei(big.NewFloat(0.00), 18),
		Data:     txData,
	})

	if err != nil {
		return nil, nil
	}

	auth.GasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2))
	auth.From = common.HexToAddress(c.Farm.FarmConfig.Wallet)
	auth.GasLimit = gas * 2
	auth.Context = context.Background()
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := c.SwapContract.SwapExactTokensForTokens(auth, amountIn, amountMin, []common.Address{
		common.HexToAddress(fromToken),
		common.HexToAddress(toToken),
	}, common.HexToAddress(wallet), big.NewInt(deadLine))
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *SwapRouter) AddLiquidity(tokenA, tokenB string, wishA, wishB, minA, minB *big.Int) (*types.Transaction, error) {
	auth, err := c.Farm.TokenBasic.CreateTx()
	if err != nil {
		return nil, err
	}
	toContract := common.HexToAddress(c.Address)
	wallet := c.Farm.FarmConfig.Wallet
	gasPrice, err := c.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	ABI, err := abi.JSON(strings.NewReader(string(contracts.SwapRouterABI)))
	if err != nil {
		return nil, err
		//log.Fatal(err)
	}
	nonce, err := c.Client.PendingNonceAt(context.Background(), common.HexToAddress(wallet))
	if err != nil {
		return nil, err
	}
	t := time.Now()
	deadLine := t.Add(time.Hour * 2).Unix()
	txData, err := ABI.Pack(
		"addLiquidity",
		common.HexToAddress(tokenA),
		common.HexToAddress(tokenB),
		wishA,
		wishB,
		minA,
		minB,
		common.HexToAddress(wallet),
		big.NewInt(deadLine),
	)
	gas, err := c.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     common.HexToAddress(wallet),
		To:       &toContract,
		GasPrice: new(big.Int).Mul(gasPrice, big.NewInt(2)),
		Value:    utils.ToWei(big.NewFloat(0.00), 18),
		Data:     txData,
	})

	if err != nil {
		return nil, err
	}

	auth.GasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2))
	auth.From = common.HexToAddress(c.Farm.FarmConfig.Wallet)
	auth.GasLimit = gas * 2
	auth.Context = context.Background()
	auth.Nonce = big.NewInt(int64(nonce))
	tx, err := c.SwapContract.AddLiquidity(
		auth,
		common.HexToAddress(tokenA),
		common.HexToAddress(tokenB),
		wishA,
		wishB,
		minA,
		minB,
		common.HexToAddress(wallet),
		big.NewInt(deadLine),
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
