package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// UtxoData utxo data mapping.
type UtxoData struct {
	Txid              string `json:"txid"`
	Vout              uint32 `json:"vout"`
	Amount            int64  `json:"amount"`
	Asset             string `json:"asset"`
	AssetBlinder      string `json:"assetblinder"`
	AssetCommitment   string `json:"assetcommitment"`
	AmountBlinder     string `json:"blinder"`
	AmountCommitment  string `json:"amountcommitment"`
	Descriptor        string `json:"descriptor"`
	ScriptsigTemplate string `json:"scriptsigTemplate"`
}

// TransactionCacheData transaction cache data mapping.
type TransactionCacheData struct {
	Hex   string     `json:"hex"`
	Utxos []UtxoData `json:"utxos"`
}

// NewTransactionCacheData returns a new TransactionCacheData struct.
func NewTransactionCacheData() *TransactionCacheData {
	return &TransactionCacheData{}
}

// InitializeTransactionCmd initialize transaction hex.
type InitializeTransactionCmd struct {
	cmd        string
	flagSet    *flag.FlagSet
	version    *uint
	locktime   *uint
	txFilePath *string
	isElements *bool
}

// NewInitializeTransactionCmd returns a new InitializeTransactionCmd struct.
func NewInitializeTransactionCmd() *InitializeTransactionCmd {
	return &InitializeTransactionCmd{}
}

// Command returns the command name.
func (cmd *InitializeTransactionCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *InitializeTransactionCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *InitializeTransactionCmd) Init() {
	cmd.cmd = "initializetransaction"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.version = cmd.flagSet.Uint("version", uint(2), "tx version")
	cmd.locktime = cmd.flagSet.Uint("locktime", uint(0), "locktime")
	cmd.txFilePath = cmd.flagSet.String("file", "", "transaction data file path")
	cmd.isElements = cmd.flagSet.Bool("elements", false, "elements mode")
}

// GetFlagSet returns the flag set for this command.
func (cmd *InitializeTransactionCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *InitializeTransactionCmd) Do(ctx context.Context) {
	var tx string
	var err error
	var handle uintptr
	if *cmd.isElements {
		handle, err = cfd.CfdGoInitializeConfidentialTransaction(
			uint32(*cmd.version),
			uint32(*cmd.locktime))
		if err == nil {
			defer cfd.CfdGoFreeTransactionHandle(handle)
			tx, err = cfd.CfdGoFinalizeTransaction(handle)
		}
	} else {
		handle, err = cfd.CfdGoInitializeTransaction(
			uint32(*cmd.version),
			uint32(*cmd.locktime))
		if err == nil {
			defer cfd.CfdGoFreeTransactionHandle(handle)
			tx, err = cfd.CfdGoFinalizeTransaction(handle)
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	if *cmd.txFilePath == "" {
		fmt.Printf("initialize transaction: %s\n", tx)
	} else {
		data := NewTransactionCacheData()
		data.Hex = tx

		_, err = os.Stat(*cmd.txFilePath)
		if err == nil {
			if err = os.Remove(*cmd.txFilePath); err != nil {
				fmt.Println(err)
				return
			}
		}

		jsonData, err := json.Marshal(*data)
		if err != nil {
			fmt.Println(err)
			return
		}

		var buf bytes.Buffer
		err = json.Indent(&buf, jsonData, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}
		indentJSON := buf.String()

		err = ioutil.WriteFile(*cmd.txFilePath, []byte(indentJSON), 666)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("initialize transaction:\n%s\n", indentJSON)
	}

}

// WriteTransactionCache write jsondata to file.
func WriteTransactionCache(path string, cache *TransactionCacheData) (jsonString string, err error) {
	if cache == nil {
		return "", errors.New("cahce is null")
	}
	jsonData, err := json.Marshal(*cache)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, jsonData, "", "  ")
	if err != nil {
		return "", err
	}
	indentJSON := buf.String()

	err = ioutil.WriteFile(path, []byte(indentJSON), 666)
	return indentJSON, err
}

// ReadTransactionCache read jsondata from file.
func ReadTransactionCache(path string) (cache *TransactionCacheData, err error) {
	_, err = os.Stat(path)
	if err != nil {
		return nil, errors.New("tx data file not found")
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data TransactionCacheData
	jsonString := strings.TrimSpace(string(bytes))
	err = json.Unmarshal([]byte(jsonString), &data)
	return &data, err
}
