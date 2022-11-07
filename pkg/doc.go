package flower

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed all-in-one.gotmpl
var aioTemplate string

func FlowDocsFromDirectories(dirs []string) (string, error) {
	flows, err := FlowsFromDirectories(dirs)
	if err != nil {
		return "", err
	}

	doc, err := allInOne(&flows)
	if err != nil {
		return "", err
	}

	return doc, err
}

func allInOne(flows *[]*Flow) (string, error) {
	tpl, err := template.New("aio").Funcs(sprig.FuncMap()).Parse(aioTemplate)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, flows)
	if err != nil {
		return "", err
	}

	return buf.String(), err
}
