package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	flower "github.com/alzabo/flower/pkg"
	"gopkg.in/yaml.v3"
)

func flowsFromYaml(yamlFile string) (flows []flower.Flow, err error) {
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
			fmt.Println("No steps found for flow", flowKey.Value)
			continue
		}

		codeBuf := new(bytes.Buffer)
		encoder := yaml.NewEncoder(codeBuf)

		if err = encoder.Encode(flowSteps); err != nil {
			return
		}
		if err = encoder.Close(); err != nil {
			return
		}

		flows = append(flows, flower.Flow{
			Name:     flowKey.Value,
			Doc:      flower.DocFromNode(flowKey),
			Code:     codeBuf.String(),
			Line:     flowKey.Line,
			FlowFile: flower.FlowFile{Path: yamlFile},
		})

		//fmt.Println("index:", i, "value:", c.Value, "tag:", c.Tag, "content:", c.Content)
	}
	return
}

func main() {
	var allFlows []flower.Flow
	dir, err := filepath.Abs("./test")

	if err != nil {
		fmt.Println(err)
	}

	yamlFiles, err := flower.FindYamlFiles(dir)

	if err != nil {
		fmt.Println("Could not find yaml files in", dir)
		return
	}

	for _, yamlFile := range yamlFiles {
		flows, err := flowsFromYaml(yamlFile)
		if err != nil {
			fmt.Println("An error occurred when reading flow data from file", yamlFile)
			continue
		}
		allFlows = append(allFlows, flows...)
	}

	fmt.Println("hi", allFlows)
}
