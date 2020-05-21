package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// BlindRawTransactionCmd append tx input.
type BlindRawTransactionCmd struct {
	cmd               string
	flagSet           *flag.FlagSet
	txFilePath        *string
	tx                *string
	blindingkeys      *string
	addresses         *string
	minimumRangeValue *int64
	exponent          *int64
	minimumBits       *int64
}

// BlindInput blinding key input
type BlindInput struct {
	txid        string
	vout        uint32
	blindingKey string
}

// NewBlindRawTransactionCmd returns a new BlindRawTransactionCmd struct.
func NewBlindRawTransactionCmd() *BlindRawTransactionCmd {
	return &BlindRawTransactionCmd{}
}

// Command returns the command name.
func (cmd *BlindRawTransactionCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *BlindRawTransactionCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *BlindRawTransactionCmd) Init() {
	cmd.cmd = "blindrawtransaction"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.blindingkeys = cmd.flagSet.String("blindingkeys", "",
		"blinding key data. format:[txid,vout,blindingKey|txid2,vout2,blindingKey2|...]")
	cmd.addresses = cmd.flagSet.String("addresses", "",
		"address data. format:[confidentialAddress1,confidentialAddress2,...]")
	cmd.minimumRangeValue = cmd.flagSet.Int64("minimumrangevalue", 1,
		"blind minimum range value")
	cmd.exponent = cmd.flagSet.Int64("exponent", 0, "blind exponent")
	cmd.minimumBits = cmd.flagSet.Int64("minimumbits", 52, "blind minimum bits")
}

// GetFlagSet returns the flag set for this command.
func (cmd *BlindRawTransactionCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *BlindRawTransactionCmd) Do(ctx context.Context) {
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

	option := cfd.NewCfdBlindTxOption()
	option.MinimumRangeValue = *cmd.minimumRangeValue
	option.Exponent = *cmd.exponent
	option.MinimumBits = *cmd.minimumBits

	inputs := []BlindInput{}
	keys := strings.Split(*cmd.blindingkeys, "|")
	for _, keyData := range keys {
		inputList := strings.Split(keyData, ",")
		if len(inputList) >= 3 {
			vout, err := strconv.Atoi(inputList[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			input := BlindInput{
				txid:        inputList[0],
				vout:        uint32(vout),
				blindingKey: inputList[2],
			}
			inputs = append(inputs, input)
		}
	}

	txinList := []cfd.CfdBlindInputData{}
	txoutList := []cfd.CfdBlindOutputData{}

	addrList := strings.Split(*cmd.addresses, ",")
	for _, addr := range addrList {
		if len(addr) > 0 {
			outData := cfd.CfdBlindOutputData{
				Index:               -1,
				ConfidentialAddress: addr,
				ConfidentialKey:     "",
			}
			txoutList = append(txoutList, outData)
		}
	}

	if data.Utxos != nil {
		for _, utxo := range data.Utxos {
			blindingKey := ""
			for _, input := range inputs {
				if utxo.Txid == input.txid && utxo.Vout == input.vout {
					blindingKey = input.blindingKey
					break
				}
			}
			blindInput := cfd.CfdBlindInputData{
				Txid:             utxo.Txid,
				Vout:             utxo.Vout,
				Asset:            utxo.Asset,
				AssetBlindFactor: utxo.AssetBlinder,
				Amount:           utxo.Amount,
				ValueBlindFactor: utxo.AmountBlinder,
				AssetBlindingKey: blindingKey,
				TokenBlindingKey: blindingKey,
			}
			txinList = append(txinList, blindInput)
		}
	}
	txHex, err := cfd.CfdGoBlindRawTransaction(tx, txinList, txoutList, &option)
	if err != nil {
		fmt.Println(err)
		return
	}
	if txHex == tx {
		fmt.Println("blinding fail.")
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
