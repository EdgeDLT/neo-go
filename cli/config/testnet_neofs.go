package config

import (
	"time"

	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/core/storage/dbconfig"
)

// TestnetNeoFS is the configuration for the NeoFS testnet.
var TestnetNeoFS = config.Config{
	ProtocolConfiguration: config.ProtocolConfiguration{
		Magic:              735783775,
		MaxTraceableBlocks: 2102400,
		InitialGASSupply:   52000000,
		TimePerBlock:       15 * time.Second,
		MemPoolSize:        50000,
		StandbyCommittee: []string{
			"0337f5f45e5be5aeae4a919d0787fcb743656560949061d5b8b05509b85ffbfd53",
			"020b86534a9a264d28b79155b0ec36d555ed0068eb1b0c4d40c35cc7d2f04759b8",
			"02c2efdc01181b0bc14fc19e0acb12281396c8c9ffe64458d621d781a1ded436b7",
			"026f9b40a73f29787ef5b289ac845bc43c64680fdd42fc170b1171d3c57213a89f",
			"0272350def90715494b857315c9b9c70181739eeec52d777424fef2891c3396cad",
			"03a8cee2d3877bcce5b4595578714d77ca2d47673150b8b9cd4e391b7c73b6bda3",
			"0215e735a657f6e23478728d1d0718d516bf50c06c2abd92ec7c00eba2bd7a2552",
		},
		ValidatorsCount:    7,
		VerifyTransactions: true,
		P2PSigExtensions:   true,
		SeedList: []string{
			"morph1.t5.fs.neo.org:50333",
			"morph2.t5.fs.neo.org:50333",
			"morph3.t5.fs.neo.org:50333",
			"morph4.t5.fs.neo.org:50333",
			"morph5.t5.fs.neo.org:50333",
			"morph6.t5.fs.neo.org:50333",
			"morph7.t5.fs.neo.org:50333",
		},
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.LevelDB,
			LevelDBOptions: dbconfig.LevelDBOptions{
				DataDirectoryPath: "./chains/testnet.neofs",
			},
		},
		P2P: config.P2P{
			Addresses:         []string{":50333"},
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MaxPeers:          100,
			MinPeers:          5,
			AttemptConnPeers:  20,
		},
		Relay: true,
		Consensus: config.Consensus{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/cn_wallet.json",
				Password: "pass",
			},
		},
		Oracle: config.OracleConfiguration{
			Enabled: false,
			AllowedContentTypes: []string{
				"application/json",
			},
		},
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{":50332"},
			},
			MaxGasInvoke:          100,
			EnableCORSWorkaround:  false,
			StartWhenSynchronized: false,
			TLSConfig: config.TLS{
				BasicService: config.BasicService{
					Enabled:   false,
					Addresses: []string{":50331"},
				},
				CertFile: "server.crt",
				KeyFile:  "server.key",
			},
		},
		P2PNotary: config.P2PNotary{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/notary_wallet.json",
				Password: "pass",
			},
		},
		Prometheus: config.BasicService{
			Enabled:   false,
			Addresses: []string{":2112"},
		},
		Pprof: config.BasicService{
			Enabled:   false,
			Addresses: []string{":2113"},
		},
	},
}
