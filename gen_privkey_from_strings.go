package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"strings"
)

type GenPrivkeyFromStringsCmd struct {
	cmd     string
	flagSet *flag.FlagSet
	text    *string
}

func NewGenPrivkeyFromStringsCmd() *GenPrivkeyFromStringsCmd {
	return &GenPrivkeyFromStringsCmd{}
}

func (cmd *GenPrivkeyFromStringsCmd) Command() string {
	return cmd.cmd
}

func (cmd *GenPrivkeyFromStringsCmd) Parse(args []string) {
	cmd.flagSet.Parse(args)
}

func (cmd *GenPrivkeyFromStringsCmd) Init() {
	cmd.cmd = "genprivkeyfromstrings"
	cmd.flagSet = flag.NewFlagSet(cmd.cmd, flag.ExitOnError)
	cmd.text = cmd.flagSet.String("text", "", "aaa|bbb|ccc")
}

func (cmd *GenPrivkeyFromStringsCmd) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

func (cmd *GenPrivkeyFromStringsCmd) Do(ctx context.Context) {
	texts := strings.Split(*cmd.text, "|")
	seed := ""
	for i, w := range texts {
		fmt.Printf("%d: '%s'\n", i, w)
		seed = seed + strings.Trim(w, " ")
	}

	h := sha256.New()
	_, err := h.Write([]byte(seed))
	if err != nil {
		fmt.Println(err)
		return
	}
	privkey := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("privkey: '%s'\n", privkey)
}
