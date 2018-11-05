package main

////////////////////////////////////////////////////////////////////////////////

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sabhiram/hist/types"

	"github.com/sabhiram/hist/emitter"
	_ "github.com/sabhiram/hist/emitter/console"
)

////////////////////////////////////////////////////////////////////////////////

const (
	cUsageStr = `usage: [history |] hist [-tag <t> [-outputs]] [-version]
If "cmd" is empty, it runs the "tag" command (see below).  Valid commands
include:

	-tag <val>	-	Set the start of a tag block with the string "<val>".
	-version 	- 	Print the version of the "hist" tool.

If the "-tag" is specified it is recorded.  If the shell history is piped
into the program, it is filtered until the last occurence of the tag "<val>"
(unless the "-tag" is specified on output as well).
`
)

////////////////////////////////////////////////////////////////////////////////

type CLIArgs struct {
	tag     string
	version bool
}

var (
	cli = CLIArgs{}
)

////////////////////////////////////////////////////////////////////////////////

func fatalOnError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func getArgFlagSet(cli *CLIArgs) *flag.FlagSet {
	fs := flag.NewFlagSet("hist", flag.ExitOnError)
	fs.StringVar(&cli.tag, "tag", "", "tag value to set in the history")
	fs.StringVar(&cli.tag, "t", "", "tag value to set in the history (short)")
	fs.BoolVar(&cli.version, "version", false, "print the version of this tool")
	fs.BoolVar(&cli.version, "v", false, "print the version of this tool (short)")
	return fs
}

func main() {
	if cli.version {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	// Check stdin
	info, err := os.Stdin.Stat()
	fatalOnError(err)

	// If the tag is set, remember it and we are done, the shell's history
	// will do the heavy lifting.
	if info.Size() == 0 {
		return
	}

	r := bufio.NewReader(os.Stdin)
	lds := []*types.LineDesc{}
	for i := 1; true; i++ {
		l, err := r.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		// TODO: Message should be annotated by the user.
		msg := fmt.Sprintf("Count = %d", i)
		lds = append([]*types.LineDesc{types.NewLineDesc(l, msg)}, lds...)
	}

	// tagv := ""
	i := 0
	for i = 0; i < len(lds); i++ {
		h := lds[i]
		t := strings.Split(strings.Trim(h.Line, " "), " ")
		arr := []string{}
		for _, v := range t {
			if len(v) > 0 {
				arr = append(arr, strings.TrimSpace(v))
			}
		}
		cmd := arr[1]
		h.Line = strings.Join(arr[1:], " ")
		if cmd == "hist" {
			h_cli := CLIArgs{}
			fs := getArgFlagSet(&h_cli)
			fs.Parse(arr[2:])

			if h_cli.tag == cli.tag {
				break
			} else if cli.tag == "" {
				break
			}
		}
	}

	// Trim till found tag.
	lds = lds[0:i]

	// Trim the most recent command and reverse the list
	ln := len(lds)
	for j := 1; j < ln/2; j++ {
		lds[j], lds[ln-j-1] = lds[ln-j-1], lds[j]
	}

	// Emit any enabled emitters.
	err = emitter.EmitEnabled(lds)
	fatalOnError(err)
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	// Allow any registered emitters to tie in their args
	fs := getArgFlagSet(&cli)
	fatalOnError(emitter.ParseArgs(fs))
	fatalOnError(fs.Parse(os.Args[1:]))
}
