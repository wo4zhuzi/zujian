package config

type RPCURL string
type NetInfo struct {
	Name        string
	RPC         string
	Router      string
	RewardToken string
	FarmAddress string
}

var NET_MAPPING = map[int]*NetInfo{
	0: &NetInfo{
		Name:        "Mdex",
		Router:      "0xED7d5F38C79115ca12fe6C0041abb22F0A06C300",
		RewardToken: "0x25d2e80cb6b86881fd7e07dd263fb79f4abe033c",
		FarmAddress: "0xFB03e11D93632D97a8981158A632Dd5986F5E909",
		RPC:         "https://http-mainnet-node.huobichain.com",
	},
	1: &NetInfo{
		Name:        "Pancake",
		Router:      "0x05ff2b0db69458a0750badebc4f9e13add608c7f",
		RewardToken: "0x0e09fabb73bd3ade0a17ecc321fd13a19e81ce82",
		FarmAddress: "0x73feaa1eE314F8c655E354234017bE2193C9E24E",
		RPC:         "https://bsc-dataseed1.binance.org",
	},
	2: &NetInfo{
		Name:        "Alpaca",
		Router:      "0x05ff2b0db69458a0750badebc4f9e13add608c7f",
		RewardToken: "0x8f0528ce5ef7b51152a59745befdd91d97091d2f",
		FarmAddress: "0xA625AB01B08ce023B2a342Dbb12a16f2C8489A8F",
		RPC:         "https://bsc-dataseed1.binance.org",
	},
}
