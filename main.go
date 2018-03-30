package main

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
	"strings"

	"github.com/sabhiram/hist/types"

	"github.com/sabhiram/hist/emitter"
	_ "github.com/sabhiram/hist/emitter/console"
)

////////////////////////////////////////////////////////////////////////////////

const (
	cUsageStr = `usage: hist [<cmd> [options]]

If "cmd" is empty, it runs the "tag" command (see below).  Valid commands 
include:

	tag <val>	-	Set the start of a tag block with the string "<val>".  If 
					"<val>" is not specified, the previous tag is closed.  On 
					close of a tag, if outputs are specified via "-<output>",
					the respective files will be emitted.

	version		- 	Print the version of the "hist" tool.
`
	cCmdEmpty   = ""
	cCmdTag     = "tag"
	cCmdVersion = "version"
)

var (
	cli = struct {
		cmd string // command to invoke for the `hist` tool
		tag string // tag value to add to history
	}{}
)

////////////////////////////////////////////////////////////////////////////////

func processTag(tag string) error {
	// If the tag is set, remember it and we are done, the shell's history
	// will do the heavy lifting.
	if len(tag) > 0 {
		return nil
	}

	// If the tag is empty, it marks the close of the last tag - bring up
	// the selector interface to choose which lines to feed to the formatter.
	ll := []*types.LineDesc{}

	// If the "-<output>" option is specified, the selected lines are passed
	// down to each appropriate <output> plugin.  For example:
	//   "-go" will invoke the `go` plugin which emits go exec statements.
	//   "-md" will invoke the `markdown` plugin to render markdown.
	//   "-console" will dump commented console versions of the scripts.
	return emitter.EmitEnabled(ll)
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	var err error
	switch cli.cmd {
	case cCmdEmpty, cCmdTag:
		err = processTag(cli.tag)
	case cCmdVersion:
		fmt.Printf(Version)
	default:
		err = fmt.Errorf("invalid command specified\n%s", cUsageStr)
	}

	emitter.DumpFactory()

	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	args, err := emitter.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}

	if len(args) > 0 {
		cli.cmd = strings.ToLower(args[0])
	} else {
		cli.cmd = cCmdTag
	}

	if len(args) > 1 {
		cli.tag = args[1]
	}
}
