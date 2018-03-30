package emitters

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"fmt"

	"github.com/sabhiram/hist/types"
)

////////////////////////////////////////////////////////////////////////////////

var (
	ErrInvalidEmitter = errors.New("invalid emitter")
)

////////////////////////////////////////////////////////////////////////////////

type EmitFn func([]*types.LineDesc) error

type EmitterDesc struct {
	name       string
	key, short string
	fn         EmitFn
}

////////////////////////////////////////////////////////////////////////////////

var factory map[string]*EmitterDesc

// RegisterEmitter is called by all emitter plugins.
func RegisterEmitter(name, key, short string, fn EmitFunc) {
	factory[key] = &EmitterDesc{
		name:  name,
		key:   key,
		short: short,
		fn:    fn,
	}
}

func EmitAll(args []string, lds []*types.LineDesc) error {
	for _, arg := range args {
		if arg[0] == '-' {
			fmt.Printf("Processing %s\n", arg[1:])
			e, ok := factory[arg[1:]]
			if !ok {
				continue
			}
			if err := e.fn(lds); err != nil {
				return err
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
