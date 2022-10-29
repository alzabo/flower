package main

import (
	"fmt"
	"path/filepath"

	flower "github.com/alzabo/flower/pkg"
)


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
		flows, err := flower.FlowsFromYaml(yamlFile)
		if err != nil {
			fmt.Println("An error occurred when reading flow data from file", yamlFile)
			continue
		}
		allFlows = append(allFlows, flows...)
	}

	fmt.Println("hi", allFlows)
}
