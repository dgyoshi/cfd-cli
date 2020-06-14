package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// GetExtkeypairFromMnemonicCmd deriveds a extended private key from seed.
type GetExtkeypairFromMnemonicCmd struct {
	cmd         string
	flagSet     *flag.FlagSet
	mnemonic    *string
	passphrase  *string
	language    *string
	networkType *string
	path        *string
}

// NewGetExtkeypairFromMnemonicCmd returns a new GetExtkeypairFromMnemonicCmd struct.
func NewGetExtkeypairFromMnemonicCmd() *GetExtkeypairFromMnemonicCmd {
	return &GetExtkeypairFromMnemonicCmd{}
}

// Command returns its command name.
func (cmd *GetExtkeypairFromMnemonicCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *GetExtkeypairFromMnemonicCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init getextkeypairfrommnemonic the command.
func (cmd *GetExtkeypairFromMnemonicCmd) Init() {
	cmd.cmd = "getextkeypairfrommnemonic"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.mnemonic = cmd.flagSet.String("mnemonic", "", "mnemonic words")
	cmd.passphrase = cmd.flagSet.String("passphrase", "", "passphrase")
	cmd.language = cmd.flagSet.String("lang", "en", "mnemonic language. (default: en) (en | jp | fr | it | es | zht | zhs)")
	cmd.networkType = cmd.flagSet.String("network", "", "mainnet | testnet | regtest")
	cmd.path = cmd.flagSet.String("path", "", "key path or paths. i.e. 1: m/44h/0h/0h/0/0 . 2: m/44h/0h/0h,m/44h/0h/1h")
}

// GetFlagSet returns the flag set for this command.
func (cmd *GetExtkeypairFromMnemonicCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *GetExtkeypairFromMnemonicCmd) Do(ctx context.Context) {

	if *cmd.mnemonic == "" {
		fmt.Println("mnemonic is required")
		return
	}

	networkType := cfd.KCfdNetworkMainnet
	switch *cmd.networkType {
	case "mainnet":
		networkType = cfd.KCfdNetworkMainnet
	case "testnet":
		networkType = cfd.KCfdNetworkTestnet
	case "regtest":
		networkType = cfd.KCfdNetworkRegtest
	default:
		fmt.Printf("network %s is unknown type.", *cmd.networkType)
	}

	mnemonicList := strings.Split(*cmd.mnemonic, " ")

	seed, _, err := cfd.CfdGoConvertMnemonicWordsToSeed(mnemonicList, *cmd.passphrase, *cmd.language)
	if err != nil {
		fmt.Println(err)
		return
	}

	baseXpriv, err := cfd.CfdGoCreateExtkeyFromSeed(seed, int(networkType), int(cfd.KCfdExtPrivkey))
	if err != nil {
		fmt.Println(err)
		return
	}

	paths := strings.Split(*cmd.path, ",")
	for _, path := range paths {
		xpriv := baseXpriv
		if path != "" {
			xpriv, err = cfd.CfdGoCreateExtkeyFromParentPath(xpriv, path, int(networkType), int(cfd.KCfdExtPrivkey))
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

		if len(path) == 0 {
			path = "m"
		}
		fmt.Printf("xpriv(%s): '%s',\nxpub (%s): '%s',\n", path, xpriv, path, xpub)
	}
}
