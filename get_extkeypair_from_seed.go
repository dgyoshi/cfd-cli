package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// GetExtkeypairFromSeedCmd deriveds a extended private key from seed.
type GetExtkeypairFromSeedCmd struct {
	cmd         string
	flagSet     *flag.FlagSet
	seed        *string
	networkType *string
	path        *string
}

// NewGetExtkeypairFromSeedCmd returns a new GetExtkeypairFromSeedCmd struct.
func NewGetExtkeypairFromSeedCmd() *GetExtkeypairFromSeedCmd {
	return &GetExtkeypairFromSeedCmd{}
}

// Command returns its command name.
func (cmd *GetExtkeypairFromSeedCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *GetExtkeypairFromSeedCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

func (cmd *GetExtkeypairFromSeedCmd) Init() {
	cmd.cmd = "getextkeypairfromseed"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.seed = cmd.flagSet.String("seed", "", "seed in hex format")
	cmd.networkType = cmd.flagSet.String("network", "", "mainnet | testnet | regtest | liquid | elementsregtest")
	cmd.path = cmd.flagSet.String("path", "", "key path. i.e. m/44h/0h/0h/0/0")
}

func (cmd *GetExtkeypairFromSeedCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

func (cmd *GetExtkeypairFromSeedCmd) Do(ctx context.Context) {

	if *cmd.seed == "" {
		fmt.Println("seed is required")
		return
	}

	networkType := cfd.KCfdNetworkMainnet
	if cmd.networkType != nil && len(*cmd.networkType) > 0 {
		switch *cmd.networkType {
		case "testnet":
			networkType = cfd.KCfdNetworkTestnet
		case "regtest":
			networkType = cfd.KCfdNetworkRegtest
		}
	}

	xpriv, err := cfd.CfdGoCreateExtkeyFromSeed(*cmd.seed, int(networkType), int(cfd.KCfdExtPrivkey))
	if err != nil {
		fmt.Println(err)
		return
	}

	if *cmd.path != "" {
		xpriv, err = cfd.CfdGoCreateExtkeyFromParentPath(xpriv, *cmd.path, int(networkType), int(cfd.KCfdExtPrivkey))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	xpub, err := cfd.CfdGoCreateExtPubkey(xpriv, int(networkType))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("xpriv: '%s'\nxpub: '%s'\n", xpriv, xpub)
}
