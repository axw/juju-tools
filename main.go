package main

import (
	"fmt"
	"os"

	"github.com/juju/cmd"
	"github.com/juju/loggo"

	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/juju"
)

var logger = loggo.GetLogger("juju.plugins.tools")

const doc = `

Juju tools is used to manage tools in the Juju controller.
This is provided for developers, and is not supported.
`

func main() {
	Main(os.Args)
}

// Main registers subcommands for the juju-local executable.
func Main(args []string) {
	ctx, err := cmd.DefaultContext()
	if err != nil {
		logger.Debugf("error: %v\n", err)
		os.Exit(2)
	}
	if err := juju.InitJujuHome(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}
	plugin := jujucmd.NewSuperCommand(cmd.SuperCommandParams{
		Name:        "tools",
		UsagePrefix: "juju",
		Doc:         doc,
		Purpose:     "manage tools in the controller",
		Log:         &cmd.Log{},
	})
	plugin.Register(envcmd.Wrap(&buildToolsCommand{}))
	plugin.Register(envcmd.Wrap(&uploadToolsCommand{}))
	plugin.Register(envcmd.Wrap(&listToolsCommand{}))
	os.Exit(cmd.Main(plugin, ctx, args[1:]))
}
