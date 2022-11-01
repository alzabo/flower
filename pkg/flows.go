package flower

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var cleanupExpr *regexp.Regexp = regexp.MustCompile(`(^#+\s*|#+$|[-=]{3,})`)
var commentCleanup *regexp.Regexp = regexp.MustCompile(`[-=]{3,}`)

type FlowFile struct {
	Path string
}
type Flow struct {
	FlowFile
	Name string
	Doc  string
	Code string
	Line int
}

func (f *FlowFile) Contents() ([]byte, error) {
	content, err := os.ReadFile(f.Path)
	return content, err
}

func parseDoc(d string) string {
	var lines []string
	for _, l := range strings.Split(d, "\n") {
		lines = append(lines, cleanupExpr.ReplaceAllLiteralString(l, ""))
	}
	return strings.Join(lines, "\n")
}

func docFromNode(n *yaml.Node) string {
	doc := n.LineComment
	if len(n.HeadComment) > 0 {
		doc = n.HeadComment
	}
	doc = parseDoc(doc)
	return doc
}

func FlowsFromYaml(yamlFile string) (flows []Flow, err error) {
	var doc yaml.Node
	var flowNodes []*yaml.Node

	content, err := os.ReadFile(yamlFile)

	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &doc)

	if err != nil {
		return
	}

	for _, node := range doc.Content {
		for i, inner := range node.Content {
			// find the node named "flows" at column 1.
			if inner.Column == 1 && inner.Value == "flows" {
				// The node after will contain flow definitions
				flowIdx := i + 1

				// can this even happen?
				if flowIdx >= len(node.Content) {
					err = errors.New(fmt.Sprint("flows key found in", yamlFile, "but flow data could not be loaded"))
					return
				}

				flowNodes = node.Content[flowIdx].Content
				break
			}
		}
	}

	// Every 2 nodes will be
	// - string: Flow key
	// - sequence: Flow steps
	for i := 0; i < len(flowNodes); i += 2 {
		flowKey := flowNodes[i]
		flowSteps := flowNodes[i+1]

		if len(flowSteps.Content) == 0 {
			// todo: debug
			//fmt.Println("No steps found for flow", flowKey.Value)
			continue
		}

		codeBuf := new(bytes.Buffer)
		encoder := yaml.NewEncoder(codeBuf)
		encoder.SetIndent(2)

		// TODO: Add key to code output, make map
		if err = encoder.Encode(flowSteps); err != nil {
			return
		}
		if err = encoder.Close(); err != nil {
			return
		}

		flows = append(flows, Flow{
			Name:     flowKey.Value,
			Doc:      docFromNode(flowKey),
			Code:     codeBuf.String(),
			Line:     flowKey.Line,
			FlowFile: FlowFile{Path: yamlFile},
		})

	}
	return
}

func FlowsFromDirectories(dirs []string) (flows []Flow, err error) {
	//dirs, err = filepath.Abs("./test")

	if err != nil {
		fmt.Println(err)
	}

	yamlFiles, err := FindYamlFiles(dirs)

	if err != nil {
		fmt.Println("Could not find yaml files in", dirs)
		return
	}

	for _, yamlFile := range yamlFiles {
		newFlows, err := FlowsFromYaml(yamlFile)
		if err != nil {
			fmt.Println("An error occurred when reading flow data from file", yamlFile)
			continue
		}
		flows = append(flows, newFlows...)
	}
	return
}
