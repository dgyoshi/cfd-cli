package main

import (
	"context"
	"flag"
	"fmt"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// EncodeDerFromSignatureCmd encode der format from signature.
type EncodeDerFromSignatureCmd struct {
	cmd          string
	flagSet      *flag.FlagSet
	sig          *string
	sighashType  *string
	anyoneCanPay *bool
}

// NewEncodeDerFromSignatureCmd returns a new EncodeDerFromSignatureCmd struct.
func NewEncodeDerFromSignatureCmd() *EncodeDerFromSignatureCmd {
	return &EncodeDerFromSignatureCmd{}
}

// Command returns the command name.
func (cmd *EncodeDerFromSignatureCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *EncodeDerFromSignatureCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *EncodeDerFromSignatureCmd) Init() {
	cmd.cmd = "encodedersignature"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.sig = cmd.flagSet.String("signature", "", "signature")
	cmd.sighashType = cmd.flagSet.String("sighashtype", "all", "signature hash type (all, none, single)")
	cmd.anyoneCanPay = cmd.flagSet.Bool("anyonecanpay", false, "enable signature hash type anyone can pay.")
}

// GetFlagSet returns the flag set for this command.
func (cmd *EncodeDerFromSignatureCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *EncodeDerFromSignatureCmd) Do(ctx context.Context) {
	if *cmd.sig == "" {
		fmt.Println("signture is required")
		return
	}

	sighashType := int(cfd.KCfdSigHashAll)
	switch *cmd.sighashType {
	case "all":
		sighashType = int(cfd.KCfdSigHashAll)
	case "none":
		sighashType = int(cfd.KCfdSigHashNone)
	case "single":
		sighashType = int(cfd.KCfdSigHashSingle)
	default:
		fmt.Printf("sighashtype %s is unknown type.", *cmd.sighashType)
		return
	}

	derSig, err := cfd.CfdGoEncodeSignatureByDer(*cmd.sig, sighashType, *cmd.anyoneCanPay)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("der encoded signature: '%s'\n", derSig)
}
