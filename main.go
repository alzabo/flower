package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

type flow struct {
	name, doc, code, file string
	line                  int
}

type flowFile struct {
	path string
}

func (f *flowFile) contents() ([]byte, error) {
	content, err := os.ReadFile(f.path)
	return content, err
}

func flowsFromYaml(yamlFile string) (flows []flow, err error) {
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
			err = encoder.Close()
			if err != nil {
				return
			}
		}

		flowDoc := flowKey.LineComment
		if len(flowKey.HeadComment) > 0 {
			flowDoc = flowKey.HeadComment
		}

		flows = append(flows, flow{
			name: flowKey.Value,
			doc:  flowDoc,
			code: codeBuf.String(),
			line: flowKey.Line,
			file: yamlFile,
		})

		//fmt.Println("index:", i, "value:", c.Value, "tag:", c.Tag, "content:", c.Content)
	}
	return
}

func findYamlFiles(dir string) ([]string, error) {
	var files []string
	contents, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, c := range contents {
		if c.IsDir() {
			nestedFiles, err := findYamlFiles(path.Join(dir, c.Name()))
			if err != nil {
				return files, err
			}

			files = append(files, nestedFiles...)
			continue
		}

		match, err := regexp.MatchString("\\.ya?ml", c.Name())

		if err != nil {
			return files, err
		}

		if !match {
			continue
		}

		files = append(files, path.Join(dir, c.Name()))
	}

	return files, nil
}

func main() {
	var allFlows []flow
	dir, err := filepath.Abs("./test")

	if err != nil {
		fmt.Println(err)
	}

	yamlFiles, err := findYamlFiles(dir)

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
