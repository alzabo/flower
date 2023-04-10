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

type FlowTemplateBinding struct {
	Flows []*Flow
	Nodes map[string]*Node
}

func FlowDocsFromDirectories(dirs []string) error {
	flows, err := FlowsFromDirectories(dirs)
	if err != nil {
		return err
	}

	// Hack: fix this later?
	//nodes := FlowGraph(flows)

	err = allInOne(FlowTemplateBinding{Flows: flows})
	if err != nil {
		return err
	}

	return err
}

func allInOne(binding FlowTemplateBinding) error {
	tpl, err := template.New("aio").Funcs(sprig.FuncMap()).Parse(aioTemplate)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, binding)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(buf.Bytes())

	return err
}
