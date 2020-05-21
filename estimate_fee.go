package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// EstimateFeeCmd fee estimation.
type EstimateFeeCmd struct {
	cmd         string
	flagSet     *flag.FlagSet
	txFilePath  *string
	tx          *string
	isElements  *bool
	feeRate     *float64
	asset       *string
	exponent    *int64
	minimumBits *int64
}

// NewEstimateFeeCmd returns a new EstimateFeeCmd struct.
func NewEstimateFeeCmd() *EstimateFeeCmd {
	return &EstimateFeeCmd{}
}

// Command returns the command name.
func (cmd *EstimateFeeCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *EstimateFeeCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *EstimateFeeCmd) Init() {
	cmd.cmd = "estimatefee"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
	cmd.feeRate = cmd.flagSet.Float64("feerate", 20.0, "fee rate. (default: 20.0)")
	cmd.asset = cmd.flagSet.String("asset", "", "fee asset")
	cmd.exponent = cmd.flagSet.Int64("exponent", 0, "blind exponent")
	cmd.minimumBits = cmd.flagSet.Int64("minimumbits", 52, "blind minimum bits")
}

// GetFlagSet returns the flag set for this command.
func (cmd *EstimateFeeCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *EstimateFeeCmd) Do(ctx context.Context) {
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

	option := cfd.NewCfdEstimateFeeOption()
	option.EffectiveFeeRate = *cmd.feeRate
	option.UseElements = *cmd.isElements
	if *cmd.isElements {
		option.FeeAsset = *cmd.asset
		option.Exponent = *cmd.exponent
		option.MinimumBits = *cmd.minimumBits
	}

	txinList := []cfd.CfdEstimateFeeInput{}

	for _, utxo := range data.Utxos {
		feeInput := cfd.CfdEstimateFeeInput{
			Utxo: cfd.CfdUtxo{
				Txid:              utxo.Txid,
				Vout:              utxo.Vout,
				Amount:            utxo.Amount,
				Asset:             utxo.Asset,
				Descriptor:        utxo.Descriptor,
				IsIssuance:        false,
				IsBlindIssuance:   false,
				IsPegin:           false,
				PeginBtcTxSize:    0,
				FedpegScript:      "",
				ScriptSigTemplate: utxo.ScriptsigTemplate,
			},
			IsIssuance:      false,
			IsBlindIssuance: false,
			IsPegin:         false,
			PeginBtcTxSize:  0,
			FedpegScript:    "",
		}
		txinList = append(txinList, feeInput)
	}
	total, txFee, inputFee, err := cfd.CfdGoEstimateFee(tx, txinList, option)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("fee = %d (tx: %d, input: %d)\n", total, txFee, inputFee)
}
