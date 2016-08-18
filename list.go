package main

import (
	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/juju/cmd/modelcmd"
	jujuversion "github.com/juju/juju/version"
	"github.com/juju/version"
)

const listToolsCommandDoc = `

Juju tools list is used to list tools in the Juju model/controller.
`

type listToolsCommand struct {
	modelcmd.ModelCommandBase
	out          cmd.Output
	versionMajor int
	series       string
	arch         string
}

// Info implements Command.Info.
func (c *listToolsCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "list",
		Purpose: "list tools to the controller",
		Doc:     listToolsCommandDoc,
	}
}

func (c *listToolsCommand) SetFlags(f *gnuflag.FlagSet) {
	c.out.AddFlags(f, "smart", cmd.DefaultFormatters)
	f.IntVar(&c.versionMajor, "major", jujuversion.Current.Major, "filter tools by major version")
	f.StringVar(&c.series, "series", "", "filter tools by series")
	f.StringVar(&c.arch, "arch", "", "filter tools by architecture")
}

// Init implements Command.Init.
func (c *listToolsCommand) Init(args []string) error {
	return nil
}

type toolsInfo struct {
	URL    string
	Size   int64
	SHA256 string
}

// Run implements Command.Run.
func (c *listToolsCommand) Run(ctx *cmd.Context) error {
	conn, err := c.NewAPIRoot()
	if err != nil {
		return errors.Annotate(err, "connecting to Juju")
	}
	defer conn.Close()
	client := conn.Client()

	result, err := client.FindTools(c.versionMajor, -1, c.series, c.arch)
	if err != nil {
		return err
	}
	if result.Error != nil {
		return result.Error
	}

	versions := make(map[version.Binary][]toolsInfo)
	for _, tools := range result.List {
		versions[tools.Version] = append(versions[tools.Version], toolsInfo{
			URL:    tools.URL,
			Size:   tools.Size,
			SHA256: tools.SHA256,
		})
	}
	return c.out.Write(ctx, versions)
}
