package config

import (
	"time"

	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/core/storage/dbconfig"
)

// MainnetNeoFS is the configuration for the NeoFS mainnet.
var MainnetNeoFS = config.Config{
	ProtocolConfiguration: config.ProtocolConfiguration{
		Magic:              91414437,
		MaxTraceableBlocks: 2102400,
		InitialGASSupply:   52000000,
		TimePerBlock:       15 * time.Second,
		MemPoolSize:        50000,
		StandbyCommittee: []string{
			"026fa34ec057d74c2fdf1a18e336d0bd597ea401a0b2ad57340d5c220d09f44086",
			"039a9db2a30942b1843db673aeb0d4fd6433f74cec1d879de6343cb9fcf7628fa4",
			"0366d255e7ce23ea6f7f1e4bedf5cbafe598705b47e6ec213ef13b2f0819e8ab33",
			"023f9cb7bbe154d529d5c719fdc39feaa831a43ae03d2a4280575b60f52fa7bc52",
			"039ba959e0ab6dc616df8b803692f1c30ba9071b76b05535eb994bf5bbc402ad5f",
			"035a2a18cddafa25ad353dea5e6730a1b9fcb4b918c4a0303c4387bb9c3b816adf",
			"031f4d9c66f2ec348832c48fd3a16dfaeb59e85f557ae1e07f6696d0375c64f97b",
		},
		ValidatorsCount:    7,
		VerifyTransactions: true,
		P2PSigExtensions:   true,
		Hardforks: map[string]uint32{
			"Aspidochelone": 3000000,
			"Basilisk":      4500000,
			"Cockatrice":    5800000,
		},
		SeedList: []string{
			"morph1.fs.neo.org:40333",
			"morph2.fs.neo.org:40333",
			"morph3.fs.neo.org:40333",
			"morph4.fs.neo.org:40333",
			"morph5.fs.neo.org:40333",
			"morph6.fs.neo.org:40333",
			"morph7.fs.neo.org:40333",
		},
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.LevelDB,
			LevelDBOptions: dbconfig.LevelDBOptions{
				DataDirectoryPath: "./chains/mainnet.neofs",
			},
		},
		P2P: config.P2P{
			Addresses:         []string{":40333"},
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
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{":40332"},
			},
			MaxGasInvoke:         15,
			EnableCORSWorkaround: false,
			TLSConfig: config.TLS{
				BasicService: config.BasicService{
					Enabled:   false,
					Addresses: []string{":40331"},
				},
				CertFile: "serv.crt",
				KeyFile:  "serv.key",
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
