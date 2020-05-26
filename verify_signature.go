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

// VerifySignatureCmd verify signature.
type VerifySignatureCmd struct {
	cmd          string
	flagSet      *flag.FlagSet
	tx           *string
	txFilePath   *string
	isElements   *bool
	txid         *string
	vout         *uint
	signature    *string
	addrType     *string
	descriptor   *string
	pubkey       *string
	redeemScript *string
	sigHashType  *string
	anyoneCanPay *bool
	amount       *uint64
	commitment   *string
}

// NewVerifySignatureCmd returns a new VerifySignatureCmd struct.
func NewVerifySignatureCmd() *VerifySignatureCmd {
	return &VerifySignatureCmd{}
}

// Command returns the command name.
func (cmd *VerifySignatureCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *VerifySignatureCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *VerifySignatureCmd) Init() {
	cmd.cmd = "verifysignature"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
	cmd.txid = cmd.flagSet.String("txid", "", "txin's txid")
	cmd.vout = cmd.flagSet.Uint("vout", 0, "txin's vout")
	cmd.signature = cmd.flagSet.String("signature", "", "txin's signature")
	cmd.descriptor = cmd.flagSet.String("descriptor", "", "txin's utxo output descriptor")
	cmd.pubkey = cmd.flagSet.String("pubkey", "", "txin's utxo pubkey (not exist descriptor)")
	cmd.redeemScript = cmd.flagSet.String("script", "", "txin's utxo redeemScript (not exist descriptor)")
	cmd.addrType = cmd.flagSet.String("addresstype", "",
		"txin's utxo addressType (p2wpkh, p2wsh, p2sh-p2wpkh, p2sh-p2wsh, p2pkh, p2sh)")
	cmd.sigHashType = cmd.flagSet.String("sighashtype", "all",
		"sighashtype (all,single,none)")
	cmd.anyoneCanPay = cmd.flagSet.Bool("anyonecanpay", false, "sighash anyonecanpay flag")
	cmd.amount = cmd.flagSet.Uint64("amount", 0, "txin's utxo amount")
	cmd.commitment = cmd.flagSet.String("commitment", "", "txin's utxo amount commitment (elements mode only)")
}

// GetFlagSet returns the flag set for this command.
func (cmd *VerifySignatureCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *VerifySignatureCmd) Do(ctx context.Context) {
	var err error
	tx := *cmd.tx
	if *cmd.tx == "" && *cmd.txFilePath != "" {
		_, err = os.Stat(*cmd.txFilePath)
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
	pubkey := *cmd.pubkey
	redeemScript := *cmd.redeemScript
	if len(*cmd.descriptor) > 0 {
		tempPubkey, tempScript, tempHashType, _, err := ParseDescriptor(*cmd.descriptor, netType)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(*cmd.addrType) == 0 {
			addrType = tempHashType
		}
		if len(pubkey) == 0 {
			pubkey = tempPubkey
		}
		if len(redeemScript) == 0 {
			redeemScript = tempScript
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

	sigHashType := -1
	anyoneCanPay := *cmd.anyoneCanPay
	signature := *cmd.signature
	if len(signature) > 130 {
		// der decode
		signature, sigHashType, anyoneCanPay, err = cfd.CfdGoDecodeSignatureFromDer(signature)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if sigHashType == -1 {
		switch *cmd.sigHashType {
		case "all":
			sigHashType = int(cfd.KCfdSigHashAll)
		case "none":
			sigHashType = int(cfd.KCfdSigHashNone)
		case "single":
			sigHashType = int(cfd.KCfdSigHashSingle)
		default:
			fmt.Printf("sighashtype %s is unknown type.", *cmd.sigHashType)
			return
		}
	}

	isVerify, err := cfd.CfdGoVerifySignature(netType, tx,
		signature, addrType, pubkey, redeemScript, *cmd.txid,
		uint32(*cmd.vout), sigHashType, anyoneCanPay,
		int64(*cmd.amount), *cmd.commitment)
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
