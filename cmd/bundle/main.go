package main

import (
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/crusttech/crust-bundle/pkg/bundle"
)

func main() {
	cfg := bundle.Configure()
	cfg.RootCommandName = "crust-bundle"
	cmd := cfg.MakeCLI(cli.Context())
	cli.HandleError(cmd.Execute())
}
