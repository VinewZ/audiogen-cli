package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vinewz/audiogen-cli/cmd/cli"
	"github.com/vinewz/audiogen-cli/cmd/utils"
)

func main() {
	utils.SupressFitzWarnings()
	if err := cli.Cmds.Run(context.Background(), os.Args); err != nil {
    fmt.Println(err)
    os.Exit(1)
	}
}

