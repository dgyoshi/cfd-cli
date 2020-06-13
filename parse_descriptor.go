package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// ParseDescriptorCmd verify signature.
type ParseDescriptorCmd struct {
	cmd        string
	flagSet    *flag.FlagSet
	descriptor *string
	nettype    *string
	childNum   *uint
}

// NewParseDescriptorCmd returns a new ParseDescriptorCmd struct.
func NewParseDescriptorCmd() *ParseDescriptorCmd {
	return &ParseDescriptorCmd{}
}

// Command returns the command name.
func (cmd *ParseDescriptorCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *ParseDescriptorCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *ParseDescriptorCmd) Init() {
	cmd.cmd = "parsedescriptor"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.descriptor = cmd.flagSet.String("descriptor", "", "txin's utxo output descriptor")
	cmd.nettype = cmd.flagSet.String("network", "mainnet", "network type (mainnet/testnet/regtest/liquidv1/liquidregtest)")
	cmd.childNum = cmd.flagSet.Uint("childnum", uint(0), "derive child number")
}

// GetFlagSet returns the flag set for this command.
func (cmd *ParseDescriptorCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *ParseDescriptorCmd) Do(ctx context.Context) {
	networkType := int(cfd.KCfdNetworkMainnet)
	switch *cmd.nettype {
	case "mainnet":
		networkType = int(cfd.KCfdNetworkMainnet)
	case "testnet":
		networkType = int(cfd.KCfdNetworkTestnet)
	case "regtest":
		networkType = int(cfd.KCfdNetworkRegtest)
	case "liquidv1":
		networkType = int(cfd.KCfdNetworkLiquidv1)
	case "liquidregtest":
		networkType = int(cfd.KCfdNetworkElementsRegtest)
	case "elementsregtest":
		networkType = int(cfd.KCfdNetworkElementsRegtest)
	default:
		fmt.Printf("nettype %s is unknown type.", *cmd.nettype)
		return
	}

	derivePath := strconv.FormatUint(uint64(*cmd.childNum), 10)
	descList, keyList, err := cfd.CfdGoParseDescriptor(*cmd.descriptor, networkType, derivePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(descList); i++ {
		fmt.Printf("[Depth:%d]\n", descList[i].Depth)
		fmt.Printf("  - LockingScript: %s\n", descList[i].LockingScript)
		if descList[i].ScriptType != int(cfd.KCfdDescriptorScriptRaw) {
			fmt.Printf("  - Address      : %s\n", descList[i].Address)
			hashType := ""
			switch descList[i].HashType {
			case int(cfd.KCfdP2pkh):
				hashType = "p2pkh"
			case int(cfd.KCfdP2sh):
				hashType = "p2sh"
			case int(cfd.KCfdP2wpkh):
				hashType = "p2wpkh"
			case int(cfd.KCfdP2wsh):
				hashType = "p2wsh"
			case int(cfd.KCfdP2shP2wpkh):
				hashType = "p2sh-p2wpkh"
			case int(cfd.KCfdP2shP2wsh):
				hashType = "p2sh-p2wsh"
			default:
				break
			}
			fmt.Printf("  - Type         : %s\n", hashType)
		}
		if (descList[i].ScriptType == int(cfd.KCfdDescriptorScriptSh)) ||
			(descList[i].ScriptType == int(cfd.KCfdDescriptorScriptWsh)) {
			fmt.Printf("  - RedeemScript : %s\n", descList[i].RedeemScript)
			scripts, err := cfd.CfdGoParseScript(descList[i].RedeemScript)
			if err == nil {
				fmt.Printf("                -> %s\n", strings.Join(scripts, " "))
			}
		}
		if descList[i].IsMultisig {
			fmt.Printf("  - requireNum   : %d\n", descList[i].ReqSigNum)
			break
		} else if descList[i].KeyType != int(cfd.KCfdDescriptorKeyNull) {
			key := ""
			if descList[i].KeyType == int(cfd.KCfdDescriptorKeyPublic) {
				key = descList[i].Pubkey
			} else if descList[i].KeyType == int(cfd.KCfdDescriptorKeyBip32) {
				key = descList[i].ExtPubkey
			} else if descList[i].KeyType == int(cfd.KCfdDescriptorKeyBip32Priv) {
				key = descList[i].ExtPrivkey
			}
			if len(key) > 0 {
				fmt.Printf("  - key          : %s\n", key)
			}
		}
	}

	if len(keyList) > 0 {
		fmt.Println("  - multisig keys:")
	}
	for i := 0; i < len(keyList); i++ {
		key := ""
		if keyList[i].KeyType == int(cfd.KCfdDescriptorKeyPublic) {
			key = keyList[i].Pubkey
		} else if keyList[i].KeyType == int(cfd.KCfdDescriptorKeyBip32) {
			key = keyList[i].ExtPubkey
		} else if keyList[i].KeyType == int(cfd.KCfdDescriptorKeyBip32Priv) {
			key = keyList[i].ExtPrivkey
		}
		if len(key) > 0 {
			fmt.Printf("    - [%d] %s\n", i, key)
		}
	}
}
