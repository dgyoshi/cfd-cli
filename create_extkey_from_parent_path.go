package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

type CreatePubkeyFromParentPathCmd struct {
	cmd         string
	flagSet     *flag.FlagSet
	xkey        *string
	path        *string
	networkType *string
}

func NewCreatePubkeyFromParentPathCmd() *CreatePubkeyFromParentPathCmd {
	return &CreatePubkeyFromParentPathCmd{}
}

func (cmd *CreatePubkeyFromParentPathCmd) Command() string {
	return cmd.cmd
}

func (cmd *CreatePubkeyFromParentPathCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

func (cmd *CreatePubkeyFromParentPathCmd) Init() {
	cmd.cmd = "createpubkeyfromparentpath"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.xkey = cmd.flagSet.String("k", "", "")
	cmd.path = cmd.flagSet.String("p", "", "")
	cmd.networkType = cmd.flagSet.String("n", "", "")
}

func (cmd *CreatePubkeyFromParentPathCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

func (cmd *CreatePubkeyFromParentPathCmd) Do(ctx context.Context) {

	networkType := cfd.KCfdNetworkMainnet
	if len(*cmd.networkType) > 0 {
		switch *cmd.networkType {
		case "mainnet":
			networkType = cfd.KCfdNetworkMainnet
		case "testnet":
			networkType = cfd.KCfdNetworkTestnet
		case "regtest":
			networkType = cfd.KCfdNetworkRegtest
		case "liquid":
			networkType = cfd.KCfdNetworkLiquidv1
		case "elementsregtest":
			networkType = cfd.KCfdNetworkElementsRegtest

		}
	}
	childKey, err := cfd.CfdGoCreateExtkeyFromParentPath(*cmd.xkey, *cmd.path, int(networkType), 1)
	if err != nil {
		panic(err)
	}

	pubkey, err := cfd.CfdGoGetPubkeyFromExtkey(childKey, int(networkType))
	if err != nil {
		panic(err)
	}

	fmt.Printf("xpub: %s\npubkey: %s\n", childKey, pubkey)
}
