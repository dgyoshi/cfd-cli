package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// SignWithPrivkeyCmd create sighash.
type SignWithPrivkeyCmd struct {
	cmd              string
	flagSet          *flag.FlagSet
	txFilePath       *string
	tx               *string
	isElements       *bool
	txid             *string
	vout             *uint
	privkey          *string
	extpriv          *string
	bip32path        *string
	addrType         *string
	amount           *int64
	amountCommitment *string
	sigHashType      *string
	anyoneCanPay     *bool
	grindR           *bool
}

// NewSignWithPrivkeyCmd returns a new SignWithPrivkeyCmd struct.
func NewSignWithPrivkeyCmd() *SignWithPrivkeyCmd {
	return &SignWithPrivkeyCmd{}
}

// Command returns the command name.
func (cmd *SignWithPrivkeyCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *SignWithPrivkeyCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *SignWithPrivkeyCmd) Init() {
	cmd.cmd = "signwithprivkey"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
	cmd.txid = cmd.flagSet.String("txid", "", "append transaction id")
	cmd.vout = cmd.flagSet.Uint("vout", uint(0), "append transaction output number")
	cmd.privkey = cmd.flagSet.String("privkey", "", "privkey")
	cmd.extpriv = cmd.flagSet.String("extpriv", "", "ext privkey")
	cmd.bip32path = cmd.flagSet.String("bip32path", "", "derive bip32 path")
	cmd.addrType = cmd.flagSet.String("addresstype", "",
		"txin's utxo addressType (p2wpkh, p2wsh, p2sh-p2wpkh, p2sh-p2wsh, p2pkh, p2sh)")
	cmd.amount = cmd.flagSet.Int64("amount", int64(0), "utxo amount")
	cmd.amountCommitment = cmd.flagSet.String("amountcommitment", "",
		"amount commitment (for blind transaction)")
	cmd.sigHashType = cmd.flagSet.String("sighashtype", "all",
		"sighashtype (all,single,none)")
	cmd.anyoneCanPay = cmd.flagSet.Bool("anyonecanpay", false, "sighash anyonecanpay flag")
	cmd.grindR = cmd.flagSet.Bool("grindr", false, "Grind-R option")
}

// GetFlagSet returns the flag set for this command.
func (cmd *SignWithPrivkeyCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *SignWithPrivkeyCmd) Do(ctx context.Context) {
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

	var privkey string
	if len(*cmd.privkey) > 0 {
		privkey = *cmd.privkey
		if len(privkey) != 64 {
			privkey, err = cfd.CfdGoGetPrivkeyFromWif(
				*cmd.privkey, int(cfd.KCfdNetworkMainnet))
			if err != nil {
				privkey, err = cfd.CfdGoGetPrivkeyFromWif(
					*cmd.privkey, int(cfd.KCfdNetworkTestnet))
			}
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		info, err := cfd.CfdGoGetExtkeyInformation(*cmd.extpriv)
		if err != nil {
			fmt.Println(err)
			return
		}
		nettype := int(cfd.KCfdNetworkTestnet)
		if info.Version == "0488ade4" {
			nettype = int(cfd.KCfdNetworkMainnet)
		}
		extpriv := *cmd.extpriv
		if len(*cmd.bip32path) > 0 {
			extpriv, err = cfd.CfdGoCreateExtkeyFromParentPath(
				extpriv, *cmd.bip32path, nettype, int(cfd.KCfdExtPrivkey))
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		privkey, _, err = cfd.CfdGoGetPrivkeyFromExtkey(extpriv, nettype)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// parameter check
	if len(*cmd.txid) != 64 {
		fmt.Println("txid size invalid.")
		return
	}
	amountCommitment := *cmd.amountCommitment
	if len(amountCommitment) > 0 && len(amountCommitment) != 66 {
		fmt.Println("amount commitment size invalid.")
		return
	}

	pubkey, err := cfd.CfdGoGetPubkeyFromPrivkey(privkey, "", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	amount := *cmd.amount
	addrType := -1
	checkPubkey, _, tempAddrType, tempAmount, tempCommitment, err := GetDescriptorInfoFromUtxoList(
		*cmd.txid, uint32(*cmd.vout), data.Utxos)
	if tempAddrType != -1 {
		if len(*cmd.addrType) == 0 {
			addrType = tempAddrType
		}
		if amount == 0 {
			amount = tempAmount
		}
		if len(amountCommitment) == 0 {
			amountCommitment = tempCommitment
		}
		if checkPubkey != pubkey {
			fmt.Printf("unmatch pubkey. %s, %s\n", checkPubkey, pubkey)
			fmt.Printf("privkey: %s\n", privkey)
		}
	}

	if addrType == -1 {
		switch *cmd.addrType {
		case "p2pkh":
			addrType = int(cfd.KCfdP2pkh)
		case "p2sh-p2wpkh":
			addrType = int(cfd.KCfdP2shP2wpkh)
		case "p2wpkh":
			addrType = int(cfd.KCfdP2wpkh)
		default:
			fmt.Printf("addresstype %s is unknown type.", *cmd.addrType)
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

	var txHex string
	if *cmd.isElements {
		txHex, err = cfd.CfdGoAddConfidentialTxSignWithPrivkey(
			tx, *cmd.txid, uint32(*cmd.vout), addrType, pubkey,
			privkey, amount, amountCommitment, sigHashType,
			*cmd.anyoneCanPay, *cmd.grindR)
	} else {
		txHex, err = cfd.CfdGoAddTxSignWithPrivkey(int(cfd.KCfdNetworkMainnet),
			tx, *cmd.txid, uint32(*cmd.vout), addrType, pubkey,
			privkey, amount, sigHashType, *cmd.anyoneCanPay, *cmd.grindR)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	if *cmd.txFilePath != "" {
		data.Hex = txHex
		_, err = WriteTransactionCache(*cmd.txFilePath, data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
