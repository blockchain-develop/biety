package config

const (
	Sync_port   = 6666
	Cons_port   = 7777
	Rpc_port    = 5555
)

var DefaultConfig = &BietyConfig{
	Genesis : &GenesisConfig{
		SeedList: []string {
			"127.0.0.1:6666",
		},
	},
}

type BietyConfig struct {
	Genesis        *GenesisConfig
}


type GenesisConfig struct {
	SeedList      []string
}
