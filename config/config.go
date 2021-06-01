package config

type SwapConfig struct {
	NetWork *NetInfo
	Wallet  string
}

type NetWorkID int

func NewSwapConfig(netInfo *NetInfo, address string) (*SwapConfig, error) {

	return &SwapConfig{
		NetWork: netInfo,
		Wallet:  address,
	}, nil
}
