package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// VerifySignTransactionCmd verify sign from transaction.
type VerifySignTransactionCmd struct {
	cmd        string
	flagSet    *flag.FlagSet
	tx         *string
	txFilePath *string
	isElements *bool
	txid       *string
	vout       *uint
	address    *string
	addrType   *string
	descriptor *string
	amount     *uint64
	commitment *string
}

// NewVerifySignTransactionCmd returns a new VerifySignTransactionCmd struct.
func NewVerifySignTransactionCmd() *VerifySignTransactionCmd {
	return &VerifySignTransactionCmd{}
}

// Command returns the command name.
func (cmd *VerifySignTransactionCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *VerifySignTransactionCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *VerifySignTransactionCmd) Init() {
	cmd.cmd = "verifysigntransaction"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
	cmd.txid = cmd.flagSet.String("txid", "", "txin's txid")
	cmd.vout = cmd.flagSet.Uint("vout", 0, "txin's vout")
	cmd.descriptor = cmd.flagSet.String("descriptor", "", "txin's utxo output descriptor")
	cmd.address = cmd.flagSet.String("address", "", "txin's utxo address (not exist descriptor)")
	cmd.addrType = cmd.flagSet.String("addresstype", "", "txin's utxo addressType (p2wpkh, p2wsh, p2sh-p2wpkh, p2sh-p2wsh, p2pkh, p2sh)")
	cmd.amount = cmd.flagSet.Uint64("amount", 0, "txin's utxo amount")
	cmd.commitment = cmd.flagSet.String("commitment", "", "txin's utxo amount commitment (elements mode only)")
}

// GetFlagSet returns the flag set for this command.
func (cmd *VerifySignTransactionCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *VerifySignTransactionCmd) Do(ctx context.Context) {
	tx := *cmd.tx
	if *cmd.tx == "" && *cmd.txFilePath != "" {
		_, err := os.Stat(*cmd.txFilePath)
		if err != nil {
			fmt.Println("tx data file not found.")
			return
		}
		txcache, err := ReadTransactionCache(*cmd.txFilePath)
		if err == nil {
			tx = txcache.Hex
		} else {
			bytes, err := ioutil.ReadFile(*cmd.txFilePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			tx = strings.TrimSpace(string(bytes))
		}
	}

	if tx == "" {
		fmt.Println("tx is required")
		return
	}

	netType := int(cfd.KCfdNetworkMainnet)
	if *cmd.isElements {
		netType = int(cfd.KCfdNetworkLiquidv1)
	}

	addrType := -1
	address := *cmd.address
	if len(*cmd.descriptor) > 0 {
		_, _, tempHashType, tempAddr, err := ParseDescriptor(*cmd.descriptor, netType)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(*cmd.addrType) == 0 {
			addrType = tempHashType
		}
		if len(address) == 0 {
			address = tempAddr
		}
	}

	if addrType == -1 {
		switch *cmd.addrType {
		case "p2pkh":
			addrType = int(cfd.KCfdP2pkhAddress)
		case "p2sh":
			addrType = int(cfd.KCfdP2shAddress)
		case "p2sh-p2wpkh":
			addrType = int(cfd.KCfdP2shP2wpkhAddress)
		case "p2sh-p2wsh":
			addrType = int(cfd.KCfdP2shP2wshAddress)
		case "p2wpkh":
			addrType = int(cfd.KCfdP2wpkhAddress)
		case "p2wsh":
			addrType = int(cfd.KCfdP2wshAddress)
		default:
			fmt.Printf("addresstype %s is unknown type.", *cmd.addrType)
			return
		}
	}

	var isVerify bool
	var err error
	if *cmd.isElements {
		isVerify, err = cfd.CfdGoVerifyConfidentialTxSign(
			tx, *cmd.txid, uint32(*cmd.vout), address,
			addrType, "", int64(*cmd.amount), *cmd.commitment)
	} else {
		isVerify, err = cfd.CfdGoVerifyTxSign(netType, tx,
			*cmd.txid, uint32(*cmd.vout), address, addrType,
			"", int64(*cmd.amount), *cmd.commitment)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("outpoint: %s,%d\n", *cmd.txid, *cmd.vout)
	if isVerify {
		fmt.Println("verify: success.")
	} else {
		fmt.Println("verify: fail.")
	}
}
