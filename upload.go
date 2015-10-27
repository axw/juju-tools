package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/version"
	"launchpad.net/gnuflag"
)

const uploadToolsCommandDoc = `

Juju tools upload is used to upload tools to the Juju controller.
`

const toolsPrefix = "juju-"
const toolsSuffix = ".tgz"

type uploadToolsCommand struct {
	envcmd.EnvCommandBase
	archives []string
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

// Init implements Command.Init.
func (c *uploadToolsCommand) Init(args []string) error {
	if len(args) == 0 {
		return errors.New("specify one or more tools archives to upload")
	}
	c.archives = args
	return nil
}

// Run implements Command.Run.
func (c *uploadToolsCommand) Run(ctx *cmd.Context) error {
	client, err := c.NewAPIClient()
	if err != nil {
		return err
	}
	defer client.Close()

	versions := make([]version.Binary, len(c.archives))
	for i, archive := range c.archives {
		basename := filepath.Base(archive)
		if !strings.HasPrefix(basename, toolsPrefix) || !strings.HasSuffix(basename, toolsSuffix) {
			return errors.NotValidf("tools archive %q", archive)
		}
		versionString := basename[len(toolsPrefix) : len(basename)-len(toolsSuffix)]
		binary, err := version.ParseBinary(versionString)
		if err != nil {
			return errors.NotValidf("tools archive %q", archive)
		}
		versions[i] = binary
	}

	for i, archive := range c.archives {
		ctx.Infof("uploading %q", archive)
		r, err := os.Open(archive)
		if err != nil {
			return err
		}
		_, err = client.UploadTools(r, versions[i])
		r.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
