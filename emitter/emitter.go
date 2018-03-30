package emitter

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"errors"
	"flag"
	"fmt"

	"github.com/sabhiram/hist/types"
)

////////////////////////////////////////////////////////////////////////////////

var (
	ErrInvalidEmitter = errors.New("invalid emitter")
)

////////////////////////////////////////////////////////////////////////////////

// EmitFunc marks the signature of the emit contract.
type EmitFunc func([]*types.LineDesc) error

// Emitter contains the registration information for all emitters that register
// with the factory.
type Emitter struct {
	name       string
	key, short string
	emit       EmitFunc
	enabled    bool
}

////////////////////////////////////////////////////////////////////////////////

// factory keeps track of all registered emitters.
var factory = map[string]*Emitter{}

// RegisterEmitter is called by all emitter plugins.
func RegisterEmitter(name, key, short string, emit EmitFunc, en bool) {
	factory[key] = &Emitter{
		name:    name,
		key:     key,
		short:   short,
		emit:    emit,
		enabled: en,
	}
}

// ParseArgs parses all args for any and all registered emitters in the factory.
// Returns a slice of stings that were not parsed by any of the registered
// parsers.
func ParseArgs(args []string) ([]string, error) {
	fs := flag.NewFlagSet("emitter-parser", flag.ContinueOnError)
	fs.SetOutput(bytes.NewBuffer([]byte{})) // Mute output from flagset

	for _, e := range factory {
		usage := fmt.Sprintf("enable output to %s", e.name)
		fs.BoolVar(&e.enabled, e.key, false, usage)
		fs.BoolVar(&e.enabled, e.short, false, usage+" (short)")
	}

	fs.Parse(args)
	return fs.Args(), nil
}

// EmitEnabled emits all enabled and registered output emitters by processing
// `ll`.
func EmitEnabled(ll []*types.LineDesc) error {
	var err error
	for _, e := range factory {
		if e.enabled {
			if err = e.emit(ll); err != nil {
				break
			}
		}
	}
	return err
}
