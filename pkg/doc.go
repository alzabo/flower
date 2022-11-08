package flower

import (
	"bytes"
	_ "embed"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed all-in-one.gotmpl
var aioTemplate string

func FlowDocsFromDirectories(dirs []string) error {
	flows, err := FlowsFromDirectories(dirs)
	if err != nil {
		return err
	}

	err = allInOne(&flows)
	if err != nil {
		return err
	}

	return err
}

func allInOne(flows *[]*Flow) error {
	tpl, err := template.New("aio").Funcs(sprig.FuncMap()).Parse(aioTemplate)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, flows)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(buf.Bytes())

	return err
}
