package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nspcc-dev/neo-go/cli/cmdargs"
	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"
)

// NewCommands returns 'config' command.
func NewCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "config",
			Usage: "NeoGo node configuration management",
			Subcommands: []cli.Command{
				{
					Name:      "generate",
					Usage:     "generate configuration files",
					UsageText: "neo-go config generate [--privnet | --mainnet | --testnet | --unit_testnet | --mainnet_neofs | --testnet_neofs | --docker | --all] [--config-path path]",
					Action:    configGenerate,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "privnet, p", Usage: "generate private network configuration"},
						cli.BoolFlag{Name: "mainnet, m", Usage: "generate mainnet network configuration"},
						cli.BoolFlag{Name: "testnet, t", Usage: "generate testnet network configuration"},
						cli.BoolFlag{Name: "unit_testnet", Usage: "generate unit test network configuration"},
						cli.BoolFlag{Name: "mainnet_neofs", Usage: "generate mainnet NeoFS network configuration"},
						cli.BoolFlag{Name: "testnet_neofs", Usage: "generate testnet NeoFS network configuration"},
						cli.BoolFlag{Name: "docker", Usage: "generate Docker configuration"},
						cli.BoolFlag{Name: "all", Usage: "generate all networks configurations"},
						cli.StringFlag{Name: "config-path", Usage: "path to the directory where configuration files will be generated"},
					},
				},
			},
		},
	}
}

func configGenerate(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}

	filesToGenerate := make(map[string]config.Config)

	if ctx.Bool("all") {
		filesToGenerate = map[string]config.Config{
			"protocol.mainnet.yml":               Mainnet,
			"protocol.mainnet.neofs.yml":         MainnetNeoFS,
			"protocol.testnet.yml":               Testnet,
			"protocol.testnet.neofs.yml":         TestnetNeoFS,
			"protocol.unit_testnet.yml":          UnitTestnet,
			"protocol.unit_testnet.single.yml":   UnitTestnetSingle,
			"protocol.privnet.docker.one.yml":    PrivnetDockerOne,
			"protocol.privnet.docker.two.yml":    PrivnetDockerTwo,
			"protocol.privnet.docker.three.yml":  PrivnetDockerThree,
			"protocol.privnet.docker.four.yml":   PrivnetDockerFour,
			"protocol.privnet.docker.single.yml": PrivnetDockerSingle,
			"protocol.privnet.yml":               Privnet,
		}
	} else {
		if ctx.Bool("privnet") {
			filesToGenerate["protocol.privnet.yml"] = Privnet
		}
		if ctx.Bool("mainnet") {
			filesToGenerate["protocol.mainnet.yml"] = Mainnet
		}
		if ctx.Bool("testnet") {
			filesToGenerate["protocol.testnet.yml"] = Testnet
		}
		if ctx.Bool("unit_testnet") {
			filesToGenerate["protocol.unit_testnet.yml"] = UnitTestnet
			filesToGenerate["protocol.unit_testnet.single.yml"] = UnitTestnetSingle
		}
		if ctx.Bool("mainnet_neofs") {
			filesToGenerate["protocol.mainnet.neofs.yml"] = MainnetNeoFS
		}
		if ctx.Bool("testnet_neofs") {
			filesToGenerate["protocol.testnet.neofs.yml"] = TestnetNeoFS
		}
		if ctx.Bool("docker") {
			filesToGenerate["protocol.privnet.docker.one.yml"] = PrivnetDockerOne
			filesToGenerate["protocol.privnet.docker.two.yml"] = PrivnetDockerTwo
			filesToGenerate["protocol.privnet.docker.three.yml"] = PrivnetDockerThree
			filesToGenerate["protocol.privnet.docker.four.yml"] = PrivnetDockerFour
			filesToGenerate["protocol.privnet.docker.single.yml"] = PrivnetDockerSingle
		}
	}
	configDir := ctx.String("config-path")
	if configDir == "" {
		configDir = "config"
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return cli.NewExitError(fmt.Errorf("unable to create configuration directory: %w", err), 1)
	}

	for fileName, cfg := range filesToGenerate {
		if err := writeConfigToFile(cfg, filepath.Join(configDir, fileName)); err != nil {
			return err
		}
	}
	return nil
}

func writeConfigToFile(cfg config.Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to create file %s: %w", filename, err), 1)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return cli.NewExitError(fmt.Errorf("failed to encode configuration to file %s: %w", filename, err), 1)
	}
	return nil
}
