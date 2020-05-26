package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// CreateSignatureHashCmd create sighash.
type CreateSignatureHashCmd struct {
	cmd              string
	flagSet          *flag.FlagSet
	txFilePath       *string
	tx               *string
	isElements       *bool
	txid             *string
	vout             *uint
	pubkey           *string
	redeemScript     *string
	addrType         *string
	amount           *int64
	amountCommitment *string
	sigHashType      *string
	anyoneCanPay     *bool
	disablecache     *bool
}

// NewCreateSignatureHashCmd returns a new CreateSignatureHashCmd struct.
func NewCreateSignatureHashCmd() *CreateSignatureHashCmd {
	return &CreateSignatureHashCmd{}
}

// Command returns the command name.
func (cmd *CreateSignatureHashCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *CreateSignatureHashCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *CreateSignatureHashCmd) Init() {
	cmd.cmd = "createsignaturehash"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
	cmd.txid = cmd.flagSet.String("txid", "", "append transaction id")
	cmd.vout = cmd.flagSet.Uint("vout", uint(0), "append transaction output number")
	cmd.pubkey = cmd.flagSet.String("pubkey", "", "pubkey (for pubkey hash)")
	cmd.redeemScript = cmd.flagSet.String("script", "", "redeem script (for script hash)")
	cmd.addrType = cmd.flagSet.String("addresstype", "",
		"txin's utxo addressType (p2wpkh, p2wsh, p2sh-p2wpkh, p2sh-p2wsh, p2pkh, p2sh)")
	cmd.amount = cmd.flagSet.Int64("amount", int64(0), "utxo amount")
	cmd.amountCommitment = cmd.flagSet.String("amountcommitment", "",
		"amount commitment (for blind transaction)")
	cmd.sigHashType = cmd.flagSet.String("sighashtype", "all",
		"sighashtype (all,single,none)")
	cmd.anyoneCanPay = cmd.flagSet.Bool("anyonecanpay", false, "sighash anyonecanpay flag")
	cmd.disablecache = cmd.flagSet.Bool("disablecache", false, "unuse cache flag")
}

// GetFlagSet returns the flag set for this command.
func (cmd *CreateSignatureHashCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *CreateSignatureHashCmd) Do(ctx context.Context) {
	var err error
	data := NewTransactionCacheData()

	tx := *cmd.tx
	if *cmd.tx == "" && *cmd.txFilePath != "" {
		data, err = ReadTransactionCache(*cmd.txFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		tx = data.Hex
	}
	if tx == "" {
		fmt.Println("tx is required")
		return
	}

	// parameter check
	if len(*cmd.txid) != 64 {
		fmt.Println("txid size invalid.")
		return
	}
	pubkey := *cmd.pubkey
	if len(pubkey) > 0 && len(pubkey) != 66 {
		fmt.Println("asset size invalid.")
		return
	}
	amountCommitment := *cmd.amountCommitment
	if len(amountCommitment) > 0 && len(amountCommitment) != 66 {
		fmt.Println("amount commitment size invalid.")
		return
	}

	amount := *cmd.amount
	redeemScript := *cmd.redeemScript
	addrType := -1
	tempPubkey, tempScript, tempAddrType, tempAmount, tempCommitment, err := GetDescriptorInfoFromUtxoList(
		*cmd.txid, uint32(*cmd.vout), data.Utxos)
	if *cmd.disablecache == false && tempAddrType != -1 {
		if len(*cmd.addrType) == 0 {
			addrType = tempAddrType
		}
		if amount == 0 {
			amount = tempAmount
		}
		if len(amountCommitment) == 0 {
			amountCommitment = tempCommitment
		}
		if len(redeemScript) == 0 {
			redeemScript = tempScript
		}
		if len(pubkey) == 0 {
			pubkey = tempPubkey
		}
	}

	if addrType == -1 {
		switch *cmd.addrType {
		case "p2pkh":
			addrType = int(cfd.KCfdP2pkh)
		case "p2sh":
			addrType = int(cfd.KCfdP2sh)
		case "p2sh-p2wpkh":
			addrType = int(cfd.KCfdP2shP2wpkh)
		case "p2sh-p2wsh":
			addrType = int(cfd.KCfdP2shP2wsh)
		case "p2wpkh":
			addrType = int(cfd.KCfdP2wpkh)
		case "p2wsh":
			addrType = int(cfd.KCfdP2wsh)
		default:
			fmt.Printf("addresstype [%s] is unknown type.", *cmd.addrType)
			return
		}
	}

	sigHashType := int(cfd.KCfdSigHashAll)
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

	var sighash string
	if *cmd.isElements {
		sighash, err = cfd.CfdGoCreateConfidentialSighash(
			tx, *cmd.txid, uint32(*cmd.vout), addrType, pubkey,
			redeemScript, amount, amountCommitment,
			sigHashType, *cmd.anyoneCanPay)
	} else {
		sighash, err = cfd.CfdGoCreateSighash(int(cfd.KCfdNetworkMainnet),
			tx, *cmd.txid, uint32(*cmd.vout), addrType, pubkey,
			redeemScript, amount, sigHashType, *cmd.anyoneCanPay)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("signature hash: %s\n", sighash)
}
