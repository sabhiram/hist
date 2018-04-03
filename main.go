package main

////////////////////////////////////////////////////////////////////////////////

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

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
into the program, it seek until the last tag in the history (unless the 
"-tag" is specified on output as well).
`
)

var (
	cli = struct {
		tag     string // tag value to add to history
		version bool   // print app version
	}{}
)

////////////////////////////////////////////////////////////////////////////////

func fatalOnError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
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

	// Read stdin.
	// TODO:
	// 1. If we find the specified tag, then we are done and we have the list
	// 	  of items we need to emit.
	// 2. If the tag is not set, we go until we find the last tag of the hist
	//	  tool in the history buffer.
	r := bufio.NewReader(os.Stdin)
	lds := []*types.LineDesc{}
	for i := 1; true; i++ {
		l, err := r.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		msg := fmt.Sprintf("Count = %d", i)
		lds = append([]*types.LineDesc{types.NewLineDesc(l, msg)}, lds...)
	}

	i := 0
	for i = 0; i < len(lds); i++ {
		h := lds[i]
		t := strings.Split(strings.Trim(h.Line, " "), " ")
		arr := []string{}
		for _, v := range t {
			if len(v) > 0 {
				arr = append(arr, v)
			}
		}
		cmd := arr[1]
		h.Line = strings.Join(arr[1:], " ")
		if cmd == "hist" {
			fmt.Printf("Trying cmd = %s\n", h.Line)
			fi := strings.Index(h.Line, "-tag")
			if fi < 0 {
				fi = strings.Index(h.Line, "-t")
			} else {
				fi += 3
			}
			if fi < 0 {
				fmt.Printf("No tag found, but hist command\n")
				break
			} else {
				fi += 1
			}
			if h.Line[fi] == ' ' || h.Line[fi] == '=' {
				fi += 1
			}
			items := strings.Split(h.Line[fi:], " ")
			tv := items[0]
			fmt.Printf("Found tag set: %s\n", tv)
		}

		fmt.Printf("%d ==> %s\n", i, h.Line)
		<-time.After(1 * time.Second)
	}

	// Trim
	lds = lds[0:i]

	// If the "-<output>" option is specified, the selected lines are passed
	// down to each appropriate <output> plugin.  For example:
	//   "-go" will invoke the `go` plugin which emits go exec statements.
	//   "-md" will invoke the `markdown` plugin to render markdown.
	//   "-console" will dump commented console versions of the scripts.
	err = emitter.EmitEnabled(lds)
	fatalOnError(err)
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	fs := flag.NewFlagSet("hist", flag.ExitOnError)
	fs.StringVar(&cli.tag, "tag", "", "tag value to set in the history")
	fs.StringVar(&cli.tag, "t", "", "tag value to set in the history (short)")
	fs.BoolVar(&cli.version, "version", false, "print the version of this tool")
	fs.BoolVar(&cli.version, "v", false, "print the version of this tool (short)")

	// Allow any registered emitters to tie in their args
	fatalOnError(emitter.ParseArgs(fs))

	// Parse args.
	fatalOnError(fs.Parse(os.Args[1:]))
}
