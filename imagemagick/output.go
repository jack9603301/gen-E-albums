package imagemagick

import "fmt"

type OutputNode struct {
	filename string
}

func (node *OutputNode) Typeid() uint {
	return NODE_TYPE_OUTPUT
}

func (node *OutputNode) GetArgs() string {
	args := fmt.Sprintf("%s", node.filename)

	return args
}
