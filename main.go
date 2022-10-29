package main

import (
	"fmt"

	flower "github.com/alzabo/flower/pkg"
)

func main() {
	allFlows, _ := flower.FlowsFromDirectory("./test")

	fmt.Println("hi", allFlows)
}
