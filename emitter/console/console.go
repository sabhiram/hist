package console

import (
	"fmt"

	"github.com/sabhiram/hist/emitter"
	"github.com/sabhiram/hist/types"
)

func consoleEmit(ll []*types.LineDesc) error {
	fmt.Printf("Running console emit\n")
	return nil
}

func init() {
	fmt.Printf("Registering the console emitter.\n")
	emitter.RegisterEmitter("console", "console", "co", consoleEmit, false)
}
