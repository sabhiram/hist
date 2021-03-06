package console

////////////////////////////////////////////////////////////////////////////////

import (
	"os"
	"text/template"

	"github.com/sabhiram/hist/emitter"
	"github.com/sabhiram/hist/types"
)

////////////////////////////////////////////////////////////////////////////////

func consoleEmit(ll []*types.LineDesc) error {
	t, err := template.New("console").Parse(`{{ range $i, $l := .}}
# {{$l.Comment}}
{{$l.Line}}
{{end}}`)
	if err != nil {
		return err
	}

	return t.Execute(os.Stdout, ll)
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	emitter.RegisterEmitter("console", "console", "co", consoleEmit, false)
}
