package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	cfd "github.com/cryptogarageinc/cfd-go"
)

// Command represent a command that can be invoked.
type Command interface {
	Command() string
	GetFlagSet() *flag.FlagSet
	Init()
	Do(context.Context)
}

var commandMap map[string]Command

func init() {
	ret := cfd.CfdInitialize()
	if ret == (int)(cfd.KCfdIllegalArgumentError) {
		panic("Fail Initialize CFD")
	}

	commandMap = make(map[string]Command)

	for _, cmd := range [...]Command{
		NewGetPubkeyFromPrivkeyCmd(),
		NewGetExtkeypairFromSeedCmd(),
		NewGenPrivkeyFromStringsCmd(),
	} {
		cmd.Init()
		commandMap[cmd.Command()] = cmd
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Need to specify a command. Available commands are:")

		for name := range commandMap {
			fmt.Println(name)
		}

		return
	}

	cmdName := os.Args[1]

	cmd, ok := commandMap[cmdName]

	if !ok {
		fmt.Println("Unknown command ", cmdName)
		return
	}

	if err := cmd.GetFlagSet().Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags %v", err)
	}

	ctx := context.Background()

	cmd.Do(ctx)
}
