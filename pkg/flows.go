package flower

import (
	"os"

	"gopkg.in/yaml.v3"
)

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

func DocFromNode(n *yaml.Node) string {
	doc := n.LineComment
	if len(n.HeadComment) > 0 {
		doc = n.HeadComment
	}
	return doc
}
