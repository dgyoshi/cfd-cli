package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// GetCommitmentCmd get commitment from blind factor.
type GetCommitmentCmd struct {
	cmd          string
	flagSet      *flag.FlagSet
	asset        *string
	amount       *int64
	assetBlinder *string
	blinder      *string
}

// NewGetCommitmentCmd returns a new GetCommitmentCmd struct.
func NewGetCommitmentCmd() *GetCommitmentCmd {
	return &GetCommitmentCmd{}
}

// Command returns the command name.
func (cmd *GetCommitmentCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *GetCommitmentCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *GetCommitmentCmd) Init() {
	cmd.cmd = "getcommitment"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.asset = cmd.flagSet.String("asset", "", "asset")
	cmd.amount = cmd.flagSet.Int64("amount", int64(0), "amount")
	cmd.assetBlinder = cmd.flagSet.String("assetblinder", "", "asset blind factor")
	cmd.blinder = cmd.flagSet.String("blinder", "", "amount blind factor")
}

// GetFlagSet returns the flag set for this command.
func (cmd *GetCommitmentCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *GetCommitmentCmd) Do(ctx context.Context) {
	if len(*cmd.asset) != 64 {
		fmt.Println("asset length is invalid")
		return
	}
	if len(*cmd.assetBlinder) != 64 {
		fmt.Println("asset blinder is invalid")
		return
	}
	if len(*cmd.blinder) != 64 {
		fmt.Println("blinder is invalid")
		return
	}

	assetCommitment, err := cfd.CfdGoGetAssetCommitment(
		*cmd.asset, *cmd.assetBlinder)
	if err != nil {
		fmt.Println(err)
		return
	}
	amountCommitment, err := cfd.CfdGoGetAmountCommitment(
		*cmd.amount, assetCommitment, *cmd.blinder)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("assetCommitment : %s\n", assetCommitment)
	fmt.Printf("amountCommitment: %s\n", amountCommitment)
}
