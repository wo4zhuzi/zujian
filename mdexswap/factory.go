package mdexswap

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wo4zhuzi/zujian/mdexswap/contracts"
	"github.com/wo4zhuzi/zujian/utils"
)

type SwapFactory struct {
	Client       *ethclient.Client
	SwapContract *contracts.Factory
	Address      string
	Farm         *MdexFarm
}

func NewSwapFactory(address string, client *ethclient.Client, farm *MdexFarm) (*SwapFactory, error) {
	swapFactory, err := contracts.NewFactory(common.HexToAddress(address), client)
	if err != nil {
		return nil, err
	}
	return &SwapFactory{
		Address:      address,
		Client:       client,
		SwapContract: swapFactory,
		Farm:         farm,
	}, nil

}

//获取期望兑换的值
func (c *SwapFactory) WishExchange(amountIn *big.Int, fromToken, toToken string) ([]*big.Int, error) {
	return c.SwapContract.GetAmountsOut(&bind.CallOpts{}, amountIn, []common.Address{
		common.HexToAddress(fromToken),
		common.HexToAddress(toToken),
	})
}

func (c *SwapFactory) Calc(amountB *big.Int, slippage float64) *big.Int {
	minAmount := new(big.Float).SetInt(amountB)

	min := new(big.Float).Mul(minAmount, big.NewFloat(1.00-slippage))
	a := big.NewInt(0)
	res, _ := min.Int(a)
	return res
}

func (c *SwapFactory) Pair(tokenA, tokenB string) (string, error) {
	lpAddress, err := c.SwapContract.PairFor(&bind.CallOpts{}, common.HexToAddress(tokenA), common.HexToAddress(tokenB))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return lpAddress.String(), nil

}
func (c *SwapFactory) PairRatio(tokenA, tokenB string) (*big.Float, *big.Float, error) {
	reserveA, reserveB, err := c.GetReserves(tokenA, tokenB)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("reserveA = ", reserveA.String())
	fmt.Println("reserveB = ", reserveB.String())
	amountA := utils.ToWei("100", int(c.Farm.TokenAInfo.Decimals))
	amountB, err := c.SwapContract.Quote(&bind.CallOpts{}, amountA, reserveA, reserveB)
	if err != nil {
		return nil, nil, err
	}
	realAmountA := utils.ToDecimal(amountA, int(c.Farm.TokenAInfo.Decimals))
	realAmountB := utils.ToDecimal(amountB, int(c.Farm.TokenBInfo.Decimals))
	realAmountAFloat, _ := realAmountA.Float64()
	realAmountBFloat, _ := realAmountB.Float64()

	total := new(big.Float).Add(realAmountA.BigFloat(), realAmountB.BigFloat())
	log.Println("realAmountA:" + realAmountA.String())
	log.Println("realAmountB:" + amountB.String())

	log.Println("total:" + total.String())
	reserveAPairRatio := new(big.Float).Quo(big.NewFloat(realAmountAFloat), total)
	reserveBPairRatio := new(big.Float).Quo(big.NewFloat(realAmountBFloat), total)
	return reserveBPairRatio, reserveAPairRatio, nil
}

func (c *SwapFactory) TokenBPairAmount(tokenA, tokenB string, amountA *big.Int) (*big.Int, error) {
	reserveA, reserveB, err := c.GetReserves(tokenA, tokenB)
	if err != nil {
		return nil, err
	}
	minB, err := c.SwapContract.Quote(&bind.CallOpts{}, amountA, reserveA, reserveB)
	if err != nil {
		return nil, nil
	}
	return minB, nil
}

func (c *SwapFactory) TokenAPairAmount(tokenA, tokenB string, amountB *big.Int) (*big.Int, error) {
	reserveA, reserveB, err := c.GetReserves(tokenA, tokenB)
	if err != nil {
		return nil, err
	}

	minB, err := c.SwapContract.Quote(&bind.CallOpts{}, amountB, reserveB, reserveA)
	if err != nil {
		return nil, nil
	}
	return minB, nil
}

func (c *SwapFactory) GetReserves(tokenA, tokenB string) (*big.Int, *big.Int, error) {
	reserves, err := c.SwapContract.GetReserves(&bind.CallOpts{}, common.HexToAddress(tokenA), common.HexToAddress(tokenB))
	if err != nil {
		return nil, nil, err
	}
	return reserves.ReserveA, reserves.ReserveB, nil
}
