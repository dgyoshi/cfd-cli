package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// AppendTxOutCmd append tx output.
type AppendTxOutCmd struct {
	cmd             string
	flagSet         *flag.FlagSet
	txFilePath      *string
	tx              *string
	isElements      *bool
	amount          *int64
	asset           *string
	address         *string
	lockingScript   *string
	isDestroyAmount *bool
	isFee           *bool
}

// NewAppendTxOutCmd returns a new AppendTxOutCmd struct.
func NewAppendTxOutCmd() *AppendTxOutCmd {
	return &AppendTxOutCmd{}
}

// Command returns the command name.
func (cmd *AppendTxOutCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *AppendTxOutCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *AppendTxOutCmd) Init() {
	cmd.cmd = "appendtxout"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
	cmd.amount = cmd.flagSet.Int64("amount", int64(0), "utxo amount")
	cmd.asset = cmd.flagSet.String("asset", "", "utxo asset")
	cmd.address = cmd.flagSet.String("address", "", "address or confidential address")
	cmd.lockingScript = cmd.flagSet.String("lockingscript", "", "locking script")
	cmd.isDestroyAmount = cmd.flagSet.Bool("destroy", false, "destroy amount")
	cmd.isFee = cmd.flagSet.Bool("fee", false, "fee output")
}

// GetFlagSet returns the flag set for this command.
func (cmd *AppendTxOutCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *AppendTxOutCmd) Do(ctx context.Context) {
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

	// other output parameter check
	if len(*cmd.asset) > 0 && len(*cmd.asset) != 64 {
		fmt.Println("asset size invalid.")
		return
	}

	var handle uintptr
	if *cmd.isElements {
		handle, err = cfd.CfdGoInitializeConfidentialTransactionByHex(
			data.Hex)
	} else {
		handle, err = cfd.CfdGoInitializeTransactionByHex(data.Hex)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cfd.CfdGoFreeTransactionHandle(handle)

	if *cmd.isElements {
		if *cmd.isFee {
			err = cfd.CfdGoAddConfidentialTxOutputFee(handle,
				*cmd.asset, *cmd.amount)
		} else if *cmd.isDestroyAmount {
			err = cfd.CfdGoAddConfidentialTxOutputDestroyAmount(handle,
				*cmd.asset, *cmd.amount)
		} else if len(*cmd.lockingScript) > 0 {
			err = cfd.CfdGoAddConfidentialTxOutputByScript(handle,
				*cmd.asset, *cmd.amount, *cmd.lockingScript)
		} else {
			err = cfd.CfdGoAddConfidentialTxOutput(handle,
				*cmd.asset, *cmd.amount, *cmd.address)
		}
	} else {
		if len(*cmd.lockingScript) > 0 {
			err = cfd.CfdGoAddTxOutputByScript(handle,
				*cmd.amount, *cmd.lockingScript)
		} else {
			err = cfd.CfdGoAddTxOutput(handle, *cmd.amount, *cmd.address)
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	txHex, err := cfd.CfdGoFinalizeTransaction(handle)
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
	fmt.Printf("append txout:\n%s\n", txHex)
}
