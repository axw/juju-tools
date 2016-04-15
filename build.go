package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/juju/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/environs/tools"
	jujuversion "github.com/juju/juju/version"
	"github.com/juju/utils/arch"
	jujuos "github.com/juju/utils/os"
	"github.com/juju/utils/series"
	"github.com/juju/version"
	"launchpad.net/gnuflag"
)

const buildToolsCommandDoc = `

Juju tools build is used to build tools archives.
`

type buildToolsCommand struct {
	modelcmd.ModelCommandBase
	version version.Binary
	dir     string
	output  string
}

// Init implements Command.Init.
func (c *buildToolsCommand) Init(args []string) error {
	arg, err := cmd.ZeroOrOneArgs(args)
	if err != nil {
		return err
	}
	if arg == "" {
		c.version.Number = jujuversion.Current
		c.version.Arch = arch.HostArch()
		if c.version.Series == "" {
			c.version.Series = series.HostSeries()
		}
	} else {
		binary, err := version.ParseBinary(arg)
		if err != nil {
			return err
		}
		c.version = binary
	}
	return nil
}

// Info implements Command.Info.
func (c *buildToolsCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "build",
		Purpose: "build tools to the controller",
		Doc:     buildToolsCommandDoc,
	}
}

// SetFlags implements Command.SetFlags.
func (c *buildToolsCommand) SetFlags(f *gnuflag.FlagSet) {
	f.StringVar(&c.dir, "d", "", "set the output directory")
	f.StringVar(&c.output, "o", "", "set the output filename")
	f.StringVar(&c.version.Series, "s", "", "set the series")
}

// Run implements Command.Run.
func (c *buildToolsCommand) Run(ctx *cmd.Context) error {
	if c.dir != "" && c.output != "" {
		return errors.New("-d and -o cannot both be specified")
	}
	if c.output == "" {
		c.output = filepath.Join(c.dir, archiveFilename(c.version))
	}
	ctx.Infof("building: %v", c.output)

	tempdir, err := ioutil.TempDir("", "juju-tools-build")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempdir)

	if err := build(c.output, c.version, tempdir); err != nil {
		return err
	}
	if err := writeForceVersion(c.version.Number, tempdir); err != nil {
		return err
	}

	archiveFile, err := os.Create(c.output)
	if err != nil {
		return err
	}
	defer archiveFile.Close()
	if err := tools.Archive(archiveFile, tempdir); err != nil {
		return err
	}

	return nil
}

func build(filename string, version version.Binary, tempdir string) error {
	seriesOS, err := series.GetOSFromSeries(version.Series)
	if err != nil {
		return err
	}
	jujudPath := filepath.Join(tempdir, "jujud")
	if seriesOS == jujuos.Windows {
		jujudPath += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", jujudPath, "github.com/juju/juju/cmd/jujud")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	env := os.Environ()
	env = environWith(env, "GOOS", osGOOS(seriesOS))
	env = environWith(env, "GOARCH", archGOARCH(version.Arch))
	cmd.Env = env

	return cmd.Run()
}

func writeForceVersion(version version.Number, tempdir string) error {
	return ioutil.WriteFile(
		filepath.Join(tempdir, "FORCE-VERSION"),
		[]byte(version.String()),
		0666,
	)
}

func archiveFilename(v version.Binary) string {
	return fmt.Sprintf("juju-%s.tgz", v)
}

func environWith(env []string, k, v string) []string {
	prefix := k + "="
	for i, kv := range env {
		if strings.HasPrefix(kv, prefix) {
			env[i] = prefix + v
			return env
		}
	}
	return append(env, prefix+v)
}

func osGOOS(os jujuos.OSType) string {
	switch os {
	case jujuos.Ubuntu, jujuos.CentOS, jujuos.Arch:
		return "linux"
	case jujuos.Windows:
		return "windows"
	}
	panic(fmt.Sprintf("unknown OS %q", os))
}

func archGOARCH(arch string) string {
	switch arch {
	case "amd64", "i386", "arm64", "ppc64":
		return arch
	case "armhf":
		return "arm"
	case "ppc64el":
		return "ppc64"
	}
	panic(fmt.Sprintf("unknown arch %q", arch))
}
