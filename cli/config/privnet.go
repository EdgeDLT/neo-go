package config

import (
	"time"

	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/core/storage/dbconfig"
)

// Privnet is the configuration for the Neo N3 privnet.
var Privnet = config.Config{
	ProtocolConfiguration: config.ProtocolConfiguration{
		Magic:              56753,
		MaxTraceableBlocks: 200000,
		TimePerBlock:       15 * time.Second,
		MemPoolSize:        50000,
		StandbyCommittee: []string{
			"02b3622bf4017bdfe317c58aed5f4c753f206b7db896046fa7d774bbc4bf7f8dc2",
			"02103a7f7dd016558597f7960d27c516a4394fd968b9e65155eb4b013e4040406e",
			"03d90c07df63e690ce77912e10ab51acc944b66860237b608c4f8f8309e71ee699",
			"02a7bc55fe8684e0119768d104ba30795bdcc86619e864add26156723ed185cd62",
		},
		ValidatorsCount:    4,
		VerifyTransactions: true,
		P2PSigExtensions:   false,
		SeedList: []string{
			"localhost:20333",
			"localhost:20334",
			"localhost:20335",
			"localhost:20336",
		},
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		Ledger: config.Ledger{
			SkipBlockVerification: false},
		// LogPath could be set up in case you need stdout logs to some proper file.
		// LogPath: "./log/neogo.log"
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.LevelDB,
			// Other options: 'inmemory','boltdb'
			// DB type options. Uncomment those you need in case you want to switch DB type.
			LevelDBOptions: dbconfig.LevelDBOptions{
				DataDirectoryPath: "./chains/privnet",
				// BoltDBOptions:
				// FilePath: "./chains/privnet.bolt"
			},
		},
		P2P: config.P2P{
			Addresses:         []string{":20332"}, // In form of "[host]:[port][:announcedPort]"
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MaxPeers:          10,
			AttemptConnPeers:  5,
			MinPeers:          3,
		},
		Relay: true,
		Consensus: config.Consensus{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/cn_wallet.json",
				Password: "pass",
			},
		},
		P2PNotary: config.P2PNotary{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/notary_wallet.json",
				Password: "pass",
			},
		},
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{":20331"},
			},
			MaxGasInvoke:          15,
			EnableCORSWorkaround:  false,
			SessionEnabled:        true,
			SessionExpirationTime: 180, // Higher expiration time for manual requests and tests.
			TLSConfig: config.TLS{
				BasicService: config.BasicService{
					Enabled:   false,
					Addresses: []string{":20330"},
				},
				CertFile: "serv.crt",
				KeyFile:  "serv.key",
			},
		},
		Prometheus: config.BasicService{
			Enabled:   true,
			Addresses: []string{":2112"},
		},
		Pprof: config.BasicService{
			Enabled:   false,
			Addresses: []string{":2113"},
		},
	},
}
