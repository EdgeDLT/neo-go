package config

import (
	"time"

	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/core/storage/dbconfig"
)

// PrivnetDockerOne is the configuration for the first node in the Dockerized private network.
var PrivnetDockerOne = config.Config{
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
			"node_one:20333",
			"node_two:20334",
			"node_three:20335",
			"node_four:20336",
		},
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		Ledger: config.Ledger{SkipBlockVerification: false},
		// LogPath: "./log/neogo.log"
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.LevelDB, // Change DB type if required. Other options: 'inmemory','boltdb'
			LevelDBOptions: dbconfig.LevelDBOptions{
				DataDirectoryPath: "/chains/one",
			},
			// BoltDBOptions: dbconfig.BoltDBOptions{
			// FilePath: "./chains/privnet.bolt"
			// }
		},
		P2P: config.P2P{
			Addresses:         []string{":20333"}, // Standard port configuration for this environment
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MaxPeers:          10,
			AttemptConnPeers:  5,
			MinPeers:          2,
		},
		Relay: true,
		Oracle: config.OracleConfiguration{
			Enabled: false,
			AllowedContentTypes: []string{
				"application/json",
			},
			Nodes: []string{
				"http://node_one:30333",
				"http://node_two:30334",
				"http://node_three:30335",
				"http://node_four:30336",
			},
			RequestTimeout: 5 * time.Second,
			UnlockWallet: config.Wallet{
				Path:     "/wallet1.json",
				Password: "one",
			},
		},
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{":30333"},
			},
			MaxGasInvoke:         15,
			EnableCORSWorkaround: false,
			SessionEnabled:       true,
		},
		P2PNotary: config.P2PNotary{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/notary_wallet.json",
				Password: "pass",
			},
		},
		Prometheus: config.BasicService{
			Enabled:   true,
			Addresses: []string{":20001"},
		},
		Pprof: config.BasicService{
			Enabled:   false,
			Addresses: []string{":20011"},
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

// PrivnetDockerTwo is the configuration for the second node in the Dockerized private network.
var PrivnetDockerTwo = config.Config{
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
		ValidatorsCount: 4,
		SeedList: []string{
			"node_one:20333",
			"node_two:20334",
			"node_three:20335",
			"node_four:20336",
		},
		VerifyTransactions: true,
		P2PSigExtensions:   false,
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		Ledger: config.Ledger{SkipBlockVerification: false},
		// LogPath could be set up in case you need stdout logs to some proper file.
		// LogPath: "./log/neogo.log"
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.LevelDB, // other options: 'inmemory','boltdb'
			// DB type options. Uncomment those you need in case you want to switch DB type.
			LevelDBOptions: dbconfig.LevelDBOptions{
				DataDirectoryPath: "/chains/two",
			},
			// BoltDBOptions:
			// FilePath: "./chains/privnet.bolt"
		},
		P2P: config.P2P{
			Addresses:         []string{":20334"}, // in form of "[host]:[port][:announcedPort]"
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MaxPeers:          10,
			AttemptConnPeers:  5,
			MinPeers:          2,
		},
		Relay: true,
		Oracle: config.OracleConfiguration{
			Enabled: false,
			AllowedContentTypes: []string{
				"application/json",
			},
			Nodes: []string{
				"http://node_one:30333",
				"http://node_two:30334",
				"http://node_three:30335",
				"http://node_four:30336",
			},
			RequestTimeout: 5 * time.Second,
			UnlockWallet: config.Wallet{
				Path:     "/wallet2.json",
				Password: "two",
			},
		},
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{":30334"}},
			MaxGasInvoke:         15,
			EnableCORSWorkaround: false,
			SessionEnabled:       true,
		},
		P2PNotary: config.P2PNotary{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/notary_wallet.json",
				Password: "pass",
			},
		},
		Prometheus: config.BasicService{
			Enabled:   true,
			Addresses: []string{":20002"},
		},
		Pprof: config.BasicService{
			Enabled:   false,
			Addresses: []string{":20012"},
		},
		Consensus: config.Consensus{
			Enabled: true,
			UnlockWallet: config.Wallet{
				Path:     "/wallet2.json",
				Password: "two",
			},
		},
	},
}

// PrivnetDockerThree is the configuration for the third node in the Dockerized private network.
var PrivnetDockerThree = config.Config{
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
			"node_one:20333",
			"node_two:20334",
			"node_three:20335",
			"node_four:20336",
		},
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		Ledger: config.Ledger{
			SkipBlockVerification: false},
		// LogPath could be set up in case you need stdout logs to some proper file.
		// LogPath: "./log/neogo.log"
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.LevelDB, // other options: 'inmemory','boltdb'
			// DB type options. Uncomment those you need in case you want to switch DB type.
			LevelDBOptions: dbconfig.LevelDBOptions{
				DataDirectoryPath: "/chains/three",
			},
			//    BoltDBOptions:
			//      FilePath: "./chains/privnet.bolt"
		},
		P2P: config.P2P{
			Addresses:         []string{":20335"}, // in form of "[host]:[port][:announcedPort]"
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MaxPeers:          10,
			AttemptConnPeers:  5,
			MinPeers:          2,
		},
		Relay: true,
		Oracle: config.OracleConfiguration{
			Enabled: false,
			AllowedContentTypes: []string{
				"application/json",
			},
			Nodes: []string{
				"http://node_one:30333",
				"http://node_two:30334",
				"http://node_three:30335",
				"http://node_four:30336",
			},
			RequestTimeout: 5 * time.Second,
			UnlockWallet: config.Wallet{
				Path:     "/wallet3.json",
				Password: "three",
			},
		},
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{":30335"}},
			MaxGasInvoke:         15,
			EnableCORSWorkaround: false,
			SessionEnabled:       true,
		},
		P2PNotary: config.P2PNotary{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/notary_wallet.json",
				Password: "pass",
			},
		},
		Prometheus: config.BasicService{
			Enabled:   true,
			Addresses: []string{":20003"},
		},
		Pprof: config.BasicService{
			Enabled:   false,
			Addresses: []string{":20013"},
		},
		Consensus: config.Consensus{
			Enabled: true,
			UnlockWallet: config.Wallet{
				Path:     "/wallet3.json",
				Password: "three",
			},
		},
	},
}

// PrivnetDockerFour is the configuration for the fourth node in the Dockerized private network.
var PrivnetDockerFour = config.Config{
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
		ValidatorsCount: 4,
		SeedList: []string{
			"node_one:20333",
			"node_two:20334",
			"node_three:20335",
			"node_four:20336",
		},
		VerifyTransactions: true,
		P2PSigExtensions:   false,
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		Ledger: config.Ledger{
			SkipBlockVerification: false},
		// LogPath could be set up in case you need stdout logs to some proper file.
		// LogPath: "./log/neogo.log"
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.LevelDB, // other options: 'inmemory','boltdb'
			// DB type options. Uncomment those you need in case you want to switch DB type.
			LevelDBOptions: dbconfig.LevelDBOptions{
				DataDirectoryPath: "/chains/four",
			},
			//    BoltDBOptions:
			//      FilePath: "./chains/privnet.bolt"
		},
		P2P: config.P2P{
			Addresses:         []string{":20336"}, // in form of "[host]:[port][:announcedPort]"
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MaxPeers:          10,
			AttemptConnPeers:  5,
			MinPeers:          2,
		},
		Relay: true,
		Oracle: config.OracleConfiguration{
			Enabled: false,
			AllowedContentTypes: []string{
				"application/json",
			},
			Nodes: []string{
				"http://node_one:30333",
				"http://node_two:30334",
				"http://node_three:30335",
				"http://node_four:30336",
			},
			RequestTimeout: 5 * time.Second,
			UnlockWallet: config.Wallet{
				Path:     "/wallet4.json",
				Password: "four",
			},
		},
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{":30336"}},
			MaxGasInvoke:         15,
			EnableCORSWorkaround: false,
			SessionEnabled:       true,
		},
		P2PNotary: config.P2PNotary{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/notary_wallet.json",
				Password: "pass",
			},
		},
		Prometheus: config.BasicService{
			Enabled:   true,
			Addresses: []string{":20004"},
		},
		Pprof: config.BasicService{
			Enabled:   false,
			Addresses: []string{":20014"},
		},
		Consensus: config.Consensus{
			Enabled: true,
			UnlockWallet: config.Wallet{
				Path:     "/wallet4.json",
				Password: "four",
			},
		},
	},
}

// PrivnetDockerSingle is the configuration for the single node in the Dockerized private network.
var PrivnetDockerSingle = config.Config{
	ProtocolConfiguration: config.ProtocolConfiguration{
		Magic:              56753,
		MaxTraceableBlocks: 200000,
		TimePerBlock:       1 * time.Second, // TimePerBlock adjusted to 1 second
		MemPoolSize:        50000,
		StandbyCommittee: []string{
			"02b3622bf4017bdfe317c58aed5f4c753f206b7db896046fa7d774bbc4bf7f8dc2",
		},
		ValidatorsCount: 1,
		SeedList: []string{
			"node_single:20333",
		},
		VerifyTransactions: true,
		P2PSigExtensions:   false,
	},

	ApplicationConfiguration: config.ApplicationConfiguration{
		Ledger: config.Ledger{
			SkipBlockVerification: false},
		// LogPath could be set up in case you need stdout logs to some proper file.
		// LogPath: "./log/neogo.log"
		DBConfiguration: dbconfig.DBConfiguration{
			Type: dbconfig.LevelDB, // other options: 'inmemory','boltdb'
			// DB type options. Uncomment those you need in case you want to switch DB type.
			LevelDBOptions: dbconfig.LevelDBOptions{
				DataDirectoryPath: "/chains/single",
			},
			//    BoltDBOptions:
			//      FilePath: "./chains/privnet.bolt"
		},
		P2P: config.P2P{
			Addresses:         []string{":20333"}, // Standard port configuration for single-node
			DialTimeout:       3 * time.Second,
			ProtoTickInterval: 2 * time.Second,
			PingInterval:      30 * time.Second,
			PingTimeout:       90 * time.Second,
			MaxPeers:          10,
			AttemptConnPeers:  5,
			MinPeers:          0,
		},
		Relay: true,
		Oracle: config.OracleConfiguration{
			Enabled: false,
			AllowedContentTypes: []string{
				"application/json",
			},
			Nodes: []string{
				"http://node_single:30333",
			},
			RequestTimeout: 5 * time.Second,
			UnlockWallet: config.Wallet{
				Path:     "/wallet1_solo.json",
				Password: "one",
			},
		},
		RPC: config.RPC{
			BasicService: config.BasicService{
				Enabled:   true,
				Addresses: []string{":30333"}},
			MaxGasInvoke:         15,
			EnableCORSWorkaround: false,
			SessionEnabled:       true,
		},
		P2PNotary: config.P2PNotary{
			Enabled: false,
			UnlockWallet: config.Wallet{
				Path:     "/notary_wallet.json",
				Password: "pass",
			},
		},
		Prometheus: config.BasicService{
			Enabled:   true,
			Addresses: []string{":20001"},
		},
		Pprof: config.BasicService{
			Enabled:   false,
			Addresses: []string{":20011"},
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
