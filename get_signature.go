package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// GetSignatureCmd get signature from privkey.
type GetSignatureCmd struct {
	cmd       string
	flagSet   *flag.FlagSet
	sighash   *string
	privkey   *string
	extpriv   *string
	bip32path *string
	grindR    *bool
}

// NewGetSignatureCmd returns a new GetSignatureCmd struct.
func NewGetSignatureCmd() *GetSignatureCmd {
	return &GetSignatureCmd{}
}

// Command returns the command name.
func (cmd *GetSignatureCmd) Command() string {
	return cmd.cmd
}

// Parse parses the command arguments.
func (cmd *GetSignatureCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

// Init initializes the command.
func (cmd *GetSignatureCmd) Init() {
	cmd.cmd = "getsignature"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.sighash = cmd.flagSet.String("sighash", "", "signature hash")
	cmd.privkey = cmd.flagSet.String("privkey", "", "privkey")
	cmd.extpriv = cmd.flagSet.String("extpriv", "", "ext privkey")
	cmd.bip32path = cmd.flagSet.String("bip32path", "", "derive bip32 path")
	cmd.grindR = cmd.flagSet.Bool("grindr", false, "Grind-R option")
}

// GetFlagSet returns the flag set for this command.
func (cmd *GetSignatureCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

// Do performs the command action.
func (cmd *GetSignatureCmd) Do(ctx context.Context) {
	if *cmd.sighash == "" {
		fmt.Println("sighash is required")
		return
	}
	sighash := *cmd.sighash
	if sigList := strings.Split(sighash, ":"); len(sigList) > 1 {
		sighash = strings.TrimSpace(sigList[len(sigList)-1])
	}

	var privkey string
	var err error
	if len(*cmd.privkey) > 0 {
		privkey = *cmd.privkey
		if len(privkey) != 64 {
			privkey, err = cfd.CfdGoGetPrivkeyFromWif(
				*cmd.privkey, int(cfd.KCfdNetworkMainnet))
			if err != nil {
				privkey, err = cfd.CfdGoGetPrivkeyFromWif(
					*cmd.privkey, int(cfd.KCfdNetworkTestnet))
			}
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		info, err := cfd.CfdGoGetExtkeyInformation(*cmd.extpriv)
		if err != nil {
			fmt.Println(err)
			return
		}
		nettype := int(cfd.KCfdNetworkTestnet)
		if info.Version == "0488ade4" {
			nettype = int(cfd.KCfdNetworkMainnet)
		}

		extpriv := *cmd.extpriv
		if len(*cmd.bip32path) > 0 {
			extpriv, err = cfd.CfdGoCreateExtkeyFromParentPath(
				extpriv, *cmd.bip32path, nettype, int(cfd.KCfdExtPrivkey))
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		privkey, _, err = cfd.CfdGoGetPrivkeyFromExtkey(extpriv, nettype)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	signature, err := cfd.CfdGoCalculateEcSignature(sighash, privkey, "",
		int(cfd.KCfdNetworkMainnet), *cmd.grindR)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("signature: %s\n", signature)
}
