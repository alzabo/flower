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

// These patterns are stripped from the raw flow comment
var cleanupExpr *regexp.Regexp = regexp.MustCompile(`(^#+\s*|#+$|[-=]{3,})`)

// Match and capture sequences like ${foo.bar}. These render strangely
// In GitHub Markdown
var codeExpr *regexp.Regexp = regexp.MustCompile(`(\$\{[^}]*\})`)

type FlowFile struct {
	Path string
}

type Flow struct {
	FlowFile
	Name     string
	Doc      string
	Code     string
	Line     int
	YamlNode *yaml.Node
}

func (f *FlowFile) Contents() ([]byte, error) {
	content, err := os.ReadFile(f.Path)
	return content, err
}

func parseDoc(d string) string {
	var lines []string
	for _, line := range strings.Split(d, "\n") {
		line = cleanupExpr.ReplaceAllLiteralString(line, "")
		lines = append(lines, line)
	}
	doc := strings.TrimSpace(strings.Join(lines, "\n"))
	doc = codeExpr.ReplaceAllString(doc, "`$1`")
	return doc
}

func docFromNode(n *yaml.Node) string {
	doc := n.LineComment
	if len(n.HeadComment) > 0 {
		doc = n.HeadComment
	}
	doc = parseDoc(doc)
	return doc
}

func FlowsFromYaml(yamlFile string) (flows []*Flow, err error) {
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

	flowNodes, err = findKeyContent(&doc, "flows")
	if err != nil {
		return
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

		// TODO: Add key to code output, make map
		var codeStr string
		if codeStr, err = encodeYaml(flowSteps); err != nil {
			// TODO: debug message
		}

		flows = append(flows, &Flow{
			Name:     flowKey.Value,
			Doc:      docFromNode(flowKey),
			Code:     codeStr,
			Line:     flowKey.Line,
			FlowFile: FlowFile{Path: yamlFile},
			YamlNode: flowSteps,
		})
	}
	return
}

func FlowsFromDirectories(dirs []string) (flows []*Flow, err error) {
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
	FlowGraph(flows)
	return
}

func FlowGraph(flows []*Flow) map[string]*Node {
	nodeMap := make(map[string]*Node, 0)
	//var nodes := []Node

	for _, flow := range flows {
		nodeMap[flow.Name] = &Node{Flow: flow}
	}

	for _, flow := range flows {
		var steps []map[string]interface{}
		yaml.Unmarshal([]byte(flow.Code), &steps)
		for _, step := range steps {
			for k, v := range step {
				if k != "call" {
					continue
				}
				callee := nodeMap[fmt.Sprint(v)]
				caller := nodeMap[flow.Name]
				edge := Edge{Caller: caller, Callee: callee}
				caller.Out = append(caller.Out, &edge)
				callee.In = append(callee.In, &edge)
			}
		}
	}
	return nodeMap
}

func encodeYaml(node *yaml.Node) (string, error) {
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(node); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func findKeyContent(n *yaml.Node, key string) ([]*yaml.Node, error) {
	var wantedNodes []*yaml.Node
	for _, node := range n.Content {
		if wantedNodes != nil {
			break
		}
		for i, inner := range node.Content {
			// find the node at column 1
			if inner.Column == 1 && inner.Value == key {
				// The node after will contain the definitions
				keyIdx := i + 1

				// can this even happen in a well-formed file?
				if keyIdx >= len(node.Content) {
					err := errors.New(fmt.Sprint("key", key, "found, but data could not found"))
					return wantedNodes, err
				}

				wantedNodes = node.Content[keyIdx].Content
				break
			}
		}
	}

	return wantedNodes, nil
}
