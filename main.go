package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////

const (
	cUsageStr = `usage: hist [<cmd> [options]]

If "cmd" is empty, it runs the "tag" command (see below).  Valid commands 
include:

	tag <val>	-	Set the start of a tag block with the string "<val>".  If 
					"<val>" is not specified, the previous tag is closed.

	version		- 	Print the version of the "hist" tool.
`
	cCmdEmpty   = ""
	cCmdTag     = "tag"
	cCmdVersion = "version"
)

var (
	cli = struct {
		cmd  string
		args []string
	}{}
)

////////////////////////////////////////////////////////////////////////////////

func processTag(args []string) error {
	// 1. 	If the tag is set, remember it and we are done, the shell's history
	// 	 	will do the heavy lifting.
	// 2. 	If the tag is empty, it marks the close of the last tag - bring up
	// 		the selector interface to choose which lines to feed to the
	// 		formatter(s).
	// 3.   If the "-<output>" option is specified, the selected lines are
	//		passed down to each appropriate <output> plugin.  For example:
	//		"-go" will invoke the `go` plugin which emits go exec statements.
	//		"-md" will invoke the `markdown` plugin to render markdown.
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	var err error
	switch cli.cmd {
	case cCmdEmpty, cCmdTag:
		err = processTag(cli.args)
	case cCmdVersion:
		fmt.Printf(Version)
	default:
		err = fmt.Errorf("no command specified\n%s", cUsageStr)
	}

	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	flag.Parse()
	cli.args = flag.Args()

	if len(cli.args) > 0 {
		cli.cmd = strings.ToLower(cli.args[0])
		cli.args = cli.args[1:]
	} else {
		cli.cmd = cCmdTag
	}
}
