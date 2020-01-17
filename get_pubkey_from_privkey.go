package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// GetUserListCmd registers a user in the system.
type GetPubkeyFromPrivkeyCmd struct {
	cmd        string
	flagSet    *flag.FlagSet
	privkey    *string
	wif        *string
	isCompress *bool
}

// NewGetPubkeyFromPrivkeyCmd returns a new GetPubkeyFromPrivkeyCmd struct.
func NewGetPubkeyFromPrivkeyCmd() *GetPubkeyFromPrivkeyCmd {
	return &GetPubkeyFromPrivkeyCmd{}
}

// Command returns the command name.
func (cmd *GetPubkeyFromPrivkeyCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *GetPubkeyFromPrivkeyCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *GetPubkeyFromPrivkeyCmd) Init() {
	cmd.cmd = "getpubkeyfromprivkey"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.privkey = cmd.flagSet.String("privkey", "", "private key in hex format")
	cmd.wif = cmd.flagSet.String("wif", "", "private key in WIF format")
	cmd.isCompress = cmd.flagSet.Bool("comp", false, "flag wether compress public key")
}

// GetFlagSet returns the flag set for this command.
func (cmd *GetPubkeyFromPrivkeyCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *GetPubkeyFromPrivkeyCmd) Do(ctx context.Context, handle uintptr) {

	if *cmd.privkey == "" && *cmd.wif == "" {
		fmt.Println("privkey or wif is required")
		return
	}

	pubkey, err := cfd.CfdGoGetPubkeyFromPrivkey(handle, *cmd.privkey, *cmd.wif, *cmd.isCompress)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("public key: '%s'\n", pubkey)
}
