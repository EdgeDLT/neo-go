package config

import (
	"time"

	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/core/storage/dbconfig"
)

// UnitTestnet is the configuration for the unit test network with four validators.
var UnitTestnet = config.Config{
	ProtocolConfiguration: config.ProtocolConfiguration{
		Magic:              42,
		MaxTraceableBlocks: 200000,
		TimePerBlock:       15 * time.Second,
		MemPoolSize:        50000,
		StandbyCommittee: []string{
			"02b3622bf4017bdfe317c58aed5f4c753f206b7db896046fa7d774bbc4bf7f8dc2",
			"02103a7f7dd016558597f7960d27c516a4394fd968b9e65155eb4b013e4040406e",
			"03d90c07df63e690ce77912e10ab51acc944b66860237b608c4f8f8309e71ee699",
			"02a7bc55fe8684e0119768d104ba30795bdcc86619e864add26156723ed185cd62",
			"02c429b3ea1aa486cb2edfd6e99d8055c1f81f1a9206664e2c40a586d187257557",
			"02c4de32252c50fa171dbe25379e4e2d55cdc12f69e382c39f59a44573ecff2f9d",
		},
		ValidatorsCount:    4,
		VerifyTransactions: true,
		P2PSigExtensions:   true,
		Hardforks: map[string]uint32{
			"Aspidochelone": 25,
		},
		SeedList: []string{
			"127.0.0.1:20334",
			"127.0.0.1:20335",
			"127.0.0.1:20336",
		},
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		Ledger: config.Ledger{
			SkipBlockVerification: false},
		// LogPath could be set up in case you need stdout logs to some proper file.
		// LogPath: "./log/neogo.log"
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.InMemoryDB, // Utilizing in-memory database for unit testing
		},
		P2P: config.P2P{
			Addresses:         []string{":20333"}, // Standard port configuration for this unit test environment
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MaxPeers:          50,
			AttemptConnPeers:  5,
			MinPeers:          0,
		},
		Relay: true,
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{"127.0.0.1:0"}},
			MaxGasInvoke:              15,
			EnableCORSWorkaround:      false,
			SessionEnabled:            true,
			SessionExpirationTime:     2,
			MaxFindStorageResultItems: 2,
		},
		Prometheus: config.BasicService{
			Enabled:   false, // Disabled for unit tests
			Addresses: []string{":2112"},
		},
		Pprof: config.BasicService{
			Enabled:   false, // Disabled for unit tests
			Addresses: []string{":2113"},
		},
		Consensus: config.Consensus{
			Enabled: true,
			UnlockWallet: config.Wallet{
				Path:     "/wallet1.json",
				Password: "one",
			},
		},
	},
}

// UnitTestnetSingle is the configuration for the unit test network with one validator.
var UnitTestnetSingle = config.Config{
	ProtocolConfiguration: config.ProtocolConfiguration{
		Magic:              42,
		MaxTraceableBlocks: 200000,
		TimePerBlock:       100 * time.Millisecond,
		MemPoolSize:        100,
		StandbyCommittee: []string{
			"02b3622bf4017bdfe317c58aed5f4c753f206b7db896046fa7d774bbc4bf7f8dc2",
		},
		ValidatorsCount:    1,
		VerifyTransactions: true,
		P2PSigExtensions:   true,
		Hardforks: map[string]uint32{
			"Aspidochelone": 25,
		},
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		Ledger: config.Ledger{
			SkipBlockVerification: false},
		// LogPath could be set up in case you need stdout logs to some proper file.
		// LogPath: "./log/neogo.log"
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.InMemoryDB,
		},
		P2P: config.P2P{
			Addresses:         []string{":0"},
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MinPeers:          0,
			MaxPeers:          10,
			AttemptConnPeers:  5,
		},
		Relay: true,
		Consensus: config.Consensus{
			Enabled: true,
			UnlockWallet: config.Wallet{
				Path:     "../testdata/wallet1_solo.json",
				Password: "one",
			},
		},
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{"127.0.0.1:0"}},
			MaxGasInvoke:         15,
			EnableCORSWorkaround: false,
		},
		Prometheus: config.BasicService{
			Enabled:   false, // Disabled for unit tests
			Addresses: []string{":2112"},
		},
		Pprof: config.BasicService{
			Enabled:   false, // Disabled for unit tests
			Addresses: []string{":2113"},
		},
	},
}
