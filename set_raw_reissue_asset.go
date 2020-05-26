package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// SetRawReissueAssetCmd append tx input.
type SetRawReissueAssetCmd struct {
	cmd           string
	flagSet       *flag.FlagSet
	txFilePath    *string
	tx            *string
	txid          *string
	vout          *uint
	entropy       *string
	assetBlinder  *string
	amount        *int64
	address       *string
	lockingScript *string
}

// NewSetRawReissueAssetCmd returns a new SetRawReissueAssetCmd struct.
func NewSetRawReissueAssetCmd() *SetRawReissueAssetCmd {
	return &SetRawReissueAssetCmd{}
}

// Command returns the command name.
func (cmd *SetRawReissueAssetCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *SetRawReissueAssetCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *SetRawReissueAssetCmd) Init() {
	cmd.cmd = "setrawreissueasset"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.txid = cmd.flagSet.String("txid", "", "append transaction id")
	cmd.vout = cmd.flagSet.Uint("vout", uint(0), "append transaction output number")
	cmd.amount = cmd.flagSet.Int64("amount", int64(0), "issue amount")
	cmd.entropy = cmd.flagSet.String("entropy", "", "issue entropy")
	cmd.assetBlinder = cmd.flagSet.String("assetblinder",
		"", "asset blinder (blindingNonce)")
	cmd.address = cmd.flagSet.String("address", "", "issue sending address")
	cmd.lockingScript = cmd.flagSet.String("lockingscript", "", "locking script")
}

// GetFlagSet returns the flag set for this command.
func (cmd *SetRawReissueAssetCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *SetRawReissueAssetCmd) Do(ctx context.Context) {
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
	assetBlinder := *cmd.assetBlinder
	if len(*cmd.assetBlinder) != 64 {
		// search from data.Utxos
		isFind := false
		for _, utxo := range data.Utxos {
			if *cmd.txid == utxo.Txid && uint32(*cmd.vout) == utxo.Vout {
				if len(utxo.AssetBlinder) > 0 {
					isFind = true
					assetBlinder = utxo.AssetBlinder
					fmt.Printf("set assetblinder: %s\n", assetBlinder)
				}
				break
			}
		}
		if !isFind {
			fmt.Println("asset blinder size invalid.")
			return
		}
	}
	if len(*cmd.entropy) != 64 {
		fmt.Println("entropy size invalid.")
		return
	}

	asset, txHex, err := cfd.CfdGoSetRawReissueAsset(tx, *cmd.txid, uint32(*cmd.vout),
		*cmd.amount, assetBlinder, *cmd.entropy, *cmd.address, *cmd.lockingScript)
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
	fmt.Printf("reissue asset: %s\n", asset)
}
