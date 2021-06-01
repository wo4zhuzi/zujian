package mdexswap

import (
	"math/big"

	"github.com/wo4zhuzi/zujian/config"
	"github.com/wo4zhuzi/zujian/token"

	"github.com/ethereum/go-ethereum/ethclient"
)

type MdexFarm struct {
	//	Printer         *printer.Printer
	FarmConfig      *config.SwapConfig
	Client          *ethclient.Client
	FarmInfo        *PoolInfo
	TokenBasic      *token.TokenBasic
	TokenAInfo      *token.Token
	TokenBInfo      *token.Token
	RewardTokenInfo *token.Token
	LpTokenInfo     *token.Token
}

type PendingReward struct {
	Amount      *big.Int
	TokenAmount *big.Int
}

func NewMdexFarm(farmConfig *config.SwapConfig, client *ethclient.Client, tokenBasic *token.TokenBasic) *MdexFarm {
	//fmt.Println("NewMdexFarm = ", "NewMdexFarm")

	return &MdexFarm{
		FarmConfig: farmConfig,
		Client:     client,
		TokenBasic: tokenBasic,
		//Printer:    printer,
	}
}
