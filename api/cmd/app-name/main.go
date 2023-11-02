package main

import (
	"github.com/fredbi/go-api-skeleton/api/cmd/app-name/commands" // CHANGE_ME
	"github.com/fredbi/go-cli/cli"
)

func main() {
	cli.MustOrDie("executing API server",
		commands.Root().Execute(),
	)
}
