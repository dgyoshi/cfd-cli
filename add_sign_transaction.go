package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// AddSignTransactionCmd append tx input.
type AddSignTransactionCmd struct {
	cmd          string
	flagSet      *flag.FlagSet
	txFilePath   *string
	tx           *string
	isElements   *bool
	txid         *string
	vout         *uint
	signature    *string
	pubkey       *string
	redeemScript *string
	addrType     *string
	sigHashType  *string
	anyoneCanPay *bool
}

// NewAddSignTransactionCmd returns a new AddSignTransactionCmd struct.
func NewAddSignTransactionCmd() *AddSignTransactionCmd {
	return &AddSignTransactionCmd{}
}

// Command returns the command name.
func (cmd *AddSignTransactionCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *AddSignTransactionCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *AddSignTransactionCmd) Init() {
	cmd.cmd = "addsigntransaction"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
	cmd.txid = cmd.flagSet.String("txid", "", "append transaction id")
	cmd.vout = cmd.flagSet.Uint("vout", uint(0), "append transaction output number")
	cmd.signature = cmd.flagSet.String("signature", "", "signature")
	cmd.pubkey = cmd.flagSet.String("pubkey", "", "pubkey")
	cmd.redeemScript = cmd.flagSet.String("script", "", "redeem script")
	cmd.addrType = cmd.flagSet.String("addresstype", "",
		"txin's utxo addressType (p2wpkh, p2wsh, p2sh-p2wpkh, p2sh-p2wsh, p2pkh, p2sh)")
	cmd.sigHashType = cmd.flagSet.String("sighashtype", "all",
		"sighashtype (all,single,none)")
	cmd.anyoneCanPay = cmd.flagSet.Bool("anyonecanpay", false, "sighash anyonecanpay flag")
}

// GetFlagSet returns the flag set for this command.
func (cmd *AddSignTransactionCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *AddSignTransactionCmd) Do(ctx context.Context) {
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

	if *cmd.isElements == false {
		fmt.Println("bitcoin tx sign is not implements.")
		return
	}

	// other input parameter check
	if len(*cmd.txid) != 64 {
		fmt.Println("txid size invalid.")
		return
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

	addrType := -1
	redeemScript := *cmd.redeemScript
	pubkey := *cmd.pubkey
	tempPubkey, tempScript, tempAddrType, _, _, err := GetDescriptorInfoFromUtxoList(
		*cmd.txid, uint32(*cmd.vout), data.Utxos)
	if tempAddrType != -1 {
		if len(*cmd.addrType) == 0 {
			addrType = tempAddrType
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
			fmt.Printf("addresstype %s is unknown type.", *cmd.addrType)
			return
		}
	}

	isMulti := false
	isScript := true
	if addrType == int(cfd.KCfdP2pkh) || addrType == int(cfd.KCfdP2wpkh) ||
		addrType == int(cfd.KCfdP2shP2wpkh) {
		isScript = false
	} else {
		_, _, err = cfd.CfdGoGetAddressesFromMultisig(
			redeemScript, int(cfd.KCfdNetworkMainnet), int(cfd.KCfdP2wpkh))
		if err == nil {
			isMulti = true
		}
	}

	var txHex string
	if !isScript {
		signData := cfd.CfdSignParameter{
			Data:                *cmd.signature,
			IsDerEncode:         true,
			SighashType:         sigHashType,
			SighashAnyoneCanPay: *cmd.anyoneCanPay,
		}
		txHex, err = cfd.CfdGoAddConfidentialTxPubkeyHashSign(
			tx, *cmd.txid, uint32(*cmd.vout), addrType, pubkey,
			signData)
	} else if isMulti {
		sigList := strings.Split(*cmd.signature, ",")
		pubkeyList := strings.Split(*cmd.pubkey, ",")
		if len(pubkeyList) > 0 && len(sigList) != len(pubkeyList) {
			fmt.Println("pubkey count is unmatch signature count.")
			return
		}
		signList := []cfd.CfdMultisigSignData{}
		for index, signature := range sigList {
			if len(signature) > 0 {
				pubkey := ""
				if len(pubkeyList) > 0 {
					pubkey = pubkeyList[index]
				}
				data := cfd.CfdMultisigSignData{
					Signature:           signature,
					IsDerEncode:         true,
					SighashType:         sigHashType,
					SighashAnyoneCanPay: *cmd.anyoneCanPay,
					RelatedPubkey:       pubkey,
				}
				signList = append(signList, data)
			}
		}
		txHex, err = cfd.CfdGoAddConfidentialTxMultisigSign(
			tx, *cmd.txid, uint32(*cmd.vout), addrType,
			signList, redeemScript)
	} else {
		sigList := strings.Split(*cmd.signature, ",")
		signList := []cfd.CfdSignParameter{}
		for _, signature := range sigList {
			if len(signature) > 0 {
				data := cfd.CfdSignParameter{
					Data:                signature,
					IsDerEncode:         false,
					SighashType:         int(cfd.KCfdSigHashAll),
					SighashAnyoneCanPay: false,
				}
				signList = append(signList, data)
			}
		}
		txHex, err = cfd.CfdGoAddConfidentialTxScriptHashSign(
			tx, *cmd.txid, uint32(*cmd.vout), addrType,
			signList, redeemScript)
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

// GetDescriptorInfoFromUtxoList get descriptor info.
func GetDescriptorInfoFromUtxoList(txid string, vout uint32, utxoList []UtxoData) (
	pubkey, redeemScript string, hashType int, amount int64,
	amountCommitment string, err error) {
	hashType = -1
	if utxoList != nil {
		desc := ""
		for _, utxo := range utxoList {
			if utxo.Txid == txid && utxo.Vout == vout {
				amount = utxo.Amount
				amountCommitment = utxo.AmountCommitment
				desc = utxo.Descriptor
				break
			}
		}
		if len(desc) > 0 {
			pubkey, redeemScript, hashType, _, err = ParseDescriptor(
				desc, int(cfd.KCfdNetworkMainnet))
			if err != nil {
				fmt.Println(err)
				return "", "", -1, -1, "", err
			}
		}
	}
	return pubkey, redeemScript, hashType, amount, amountCommitment, nil
}

// ParseDescriptor parse descriptor
func ParseDescriptor(descriptor string, networkType int) (pubkey, redeemScript string, hashType int, address string, err error) {
	descList, _, err := cfd.CfdGoParseDescriptor(
		descriptor, networkType, "")
	if err != nil {
		fmt.Println(err)
		return "", "", -1, "", err
	}
	hashType = descList[0].HashType
	address = descList[0].Address
	redeemScript = descList[len(descList)-1].RedeemScript
	extkey := ""
	switch descList[len(descList)-1].KeyType {
	case int(cfd.KCfdDescriptorKeyPublic):
		pubkey = descList[len(descList)-1].Pubkey
	case int(cfd.KCfdDescriptorKeyBip32):
		extkey = descList[len(descList)-1].ExtPubkey
	case int(cfd.KCfdDescriptorKeyBip32Priv):
		extkey = descList[len(descList)-1].ExtPrivkey
	default:
		// do nothing
	}
	if len(extkey) > 0 {
		key, err := cfd.CfdGoGetPubkeyFromExtkey(
			extkey, int(cfd.KCfdNetworkMainnet))
		if err != nil {
			key, err = cfd.CfdGoGetPubkeyFromExtkey(
				extkey, int(cfd.KCfdNetworkTestnet))
		}
		if err != nil {
			fmt.Println(err)
			return "", "", -1, "", err
		}
		pubkey = key
	}
	return pubkey, redeemScript, hashType, address, nil
}
