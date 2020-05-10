package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// DecodeRawTransactionCmd decode transaction hex.
type DecodeRawTransactionCmd struct {
	cmd        string
	flagSet    *flag.FlagSet
	tx         *string
	txFilePath *string
	nettype    *string
	isElements *bool
}

// NewDecodeRawTransactionCmd returns a new DecodeRawTransactionCmd struct.
func NewDecodeRawTransactionCmd() *DecodeRawTransactionCmd {
	return &DecodeRawTransactionCmd{}
}

// Command returns the command name.
func (cmd *DecodeRawTransactionCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *DecodeRawTransactionCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *DecodeRawTransactionCmd) Init() {
	cmd.cmd = "decoderawtransaction"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.tx = cmd.flagSet.String("tx", "", "transaction in hex format")
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.nettype = cmd.flagSet.String("network", "mainnet", "network type (mainnet/testnet/regtest)")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
}

// GetFlagSet returns the flag set for this command.
func (cmd *DecodeRawTransactionCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *DecodeRawTransactionCmd) Do(ctx context.Context) {
	tx := *cmd.tx
	if *cmd.tx == "" && *cmd.txFilePath != "" {
		_, err := os.Stat(*cmd.txFilePath)
		if err != nil {
			fmt.Println("tx data file not found.")
			return
		}
		bytes, err := ioutil.ReadFile(*cmd.txFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		tx = strings.TrimSpace(string(bytes))
	}

	if tx == "" {
		fmt.Println("tx is required")
		return
	}

	jsonData, err := cfd.CfdGoDecodeRawTransactionJson(tx, *cmd.nettype, *cmd.isElements)
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(jsonData), "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	indentJSON := buf.String()

	fmt.Printf("decode transaction:\n%s\n", indentJSON)
}
