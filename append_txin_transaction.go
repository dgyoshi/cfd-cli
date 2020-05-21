package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// AppendTxInCmd append tx input.
type AppendTxInCmd struct {
	cmd               string
	flagSet           *flag.FlagSet
	txFilePath        *string
	tx                *string
	isElements        *bool
	txid              *string
	vout              *uint
	sequence          *uint
	amount            *int64
	asset             *string
	assetBlinder      *string
	assetCommitment   *string
	amountBlinder     *string
	amountCommitment  *string
	descriptor        *string
	scriptsigTemplate *string
}

// NewAppendTxInCmd returns a new AppendTxInCmd struct.
func NewAppendTxInCmd() *AppendTxInCmd {
	return &AppendTxInCmd{}
}

// Command returns the command name.
func (cmd *AppendTxInCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *AppendTxInCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *AppendTxInCmd) Init() {
	cmd.cmd = "appendtxin"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
	cmd.txid = cmd.flagSet.String("txid", "", "append transaction id")
	cmd.vout = cmd.flagSet.Uint("vout", uint(0), "append transaction output number")
	cmd.sequence = cmd.flagSet.Uint("sequence", uint(0xffffffff), "sequence number")
	cmd.amount = cmd.flagSet.Int64("amount", int64(0), "utxo amount")
	cmd.asset = cmd.flagSet.String("asset", "", "utxo asset")
	cmd.assetBlinder = cmd.flagSet.String("assetblinder",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"asset blinder (with confidential transaction)")
	cmd.amountBlinder = cmd.flagSet.String("blinder",
		"0000000000000000000000000000000000000000000000000000000000000000",
		"amount blinder (with confidential transaction)")
	cmd.assetCommitment = cmd.flagSet.String("assetcommitment", "",
		"asset commitment")
	cmd.amountCommitment = cmd.flagSet.String("amountcommitment", "",
		"amount commitment")
	cmd.descriptor = cmd.flagSet.String("descriptor", "", "output descriptor")
	cmd.scriptsigTemplate = cmd.flagSet.String("scriptsigTemplate", "",
		"scriptsig template (for estimate fee)")
}

// GetFlagSet returns the flag set for this command.
func (cmd *AppendTxInCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *AppendTxInCmd) Do(ctx context.Context) {
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

	// other input parameter check
	if len(*cmd.txid) != 64 {
		fmt.Println("txid size invalid.")
		return
	}
	if len(*cmd.asset) > 0 && len(*cmd.asset) != 64 {
		fmt.Println("asset size invalid.")
		return
	}
	if len(*cmd.assetBlinder) > 0 && len(*cmd.assetBlinder) != 64 {
		fmt.Println("asset blinder size invalid.")
		return
	}
	if len(*cmd.amountBlinder) > 0 && len(*cmd.amountBlinder) != 64 {
		fmt.Println("amount blinder size invalid.")
		return
	}
	if len(*cmd.assetCommitment) > 0 && len(*cmd.assetCommitment) != 66 {
		fmt.Println("asset commitment size invalid.")
		return
	}
	if len(*cmd.amountCommitment) > 0 && len(*cmd.amountCommitment) != 66 {
		fmt.Println("amount commitment size invalid.")
		return
	}
	if len(*cmd.descriptor) > 0 {
		netType := int(cfd.KCfdNetworkMainnet)
		if *cmd.isElements {
			netType = int(cfd.KCfdNetworkLiquidv1)
		}
		_, _, err = cfd.CfdGoParseDescriptor(*cmd.descriptor, netType, "")
		if err != nil {
			fmt.Println("descriptor is invalid.")
			fmt.Println(err)
			return
		}
	}

	var handle uintptr
	if *cmd.isElements {
		handle, err = cfd.CfdGoInitializeConfidentialTransactionByHex(tx)
	} else {
		handle, err = cfd.CfdGoInitializeTransactionByHex(tx)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cfd.CfdGoFreeTransactionHandle(handle)

	err = cfd.CfdGoAddTxInput(handle, *cmd.txid,
		uint32(*cmd.vout), uint32(*cmd.sequence))
	if err != nil {
		fmt.Println(err)
		return
	}
	txHex, err := cfd.CfdGoFinalizeTransaction(handle)
	if err != nil {
		fmt.Println(err)
		return
	}

	isFind := false
	_, err = cfd.CfdGoGetConfidentialTxInIndex(tx, *cmd.txid, uint32(*cmd.vout))
	if err == nil {
		// already exist. tx not update.
		isFind = true
		txHex = tx
	}

	if *cmd.txFilePath != "" {
		utxo := UtxoData{
			Txid:              *cmd.txid,
			Vout:              uint32(*cmd.vout),
			Amount:            *cmd.amount,
			Asset:             *cmd.asset,
			AssetBlinder:      *cmd.assetBlinder,
			AssetCommitment:   *cmd.assetCommitment,
			AmountBlinder:     *cmd.amountBlinder,
			AmountCommitment:  *cmd.amountCommitment,
			Descriptor:        *cmd.descriptor,
			ScriptsigTemplate: *cmd.scriptsigTemplate,
		}
		data.Hex = txHex
		isUpdate := false
		if isFind {
			// update utxo data
			for index, utxoData := range data.Utxos {
				if *cmd.txid == utxoData.Txid && uint32(*cmd.vout) == utxoData.Vout {
					data.Utxos[index] = utxo
					isUpdate = true
					break
				}
			}
		}
		if !isUpdate {
			data.Utxos = append(data.Utxos, utxo)
		}

		_, err = WriteTransactionCache(*cmd.txFilePath, data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Printf("append txin:\n%s\n", txHex)
}
