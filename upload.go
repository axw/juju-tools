package main

import (
	"github.com/juju/cmd"
	"github.com/juju/juju/cmd/envcmd"
	"launchpad.net/gnuflag"
)

const uploadToolsCommandDoc = `

Juju tools upload is used to upload tools to the Juju controller.
`

type uploadToolsCommand struct {
	envcmd.EnvCommandBase
}

// Init implements Command.Init.
func (c *uploadToolsCommand) Init(args []string) (err error) {
	return nil
}

// Info implements Command.Info.
func (c *uploadToolsCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "upload",
		Purpose: "upload tools to the controller",
		Doc:     uploadToolsCommandDoc,
	}
}

// SetFlags implements Command.SetFlags.
func (c *uploadToolsCommand) SetFlags(f *gnuflag.FlagSet) {
	c.EnvCommandBase.SetFlags(f)
}

// Run implements Command.Run.
func (c *uploadToolsCommand) Run(ctx *cmd.Context) error {
	return nil
}
