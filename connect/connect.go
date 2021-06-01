package connect

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wo4zhuzi/zujian/config"
)

func Connect(cfg *config.NetInfo) (client *ethclient.Client, err error) {
	client, err = ethclient.Dial(cfg.RPC)
	//cf = func() { client.Close() }
	return
}
