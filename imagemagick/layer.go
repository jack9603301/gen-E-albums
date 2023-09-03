package imagemagick

import "fmt"

type LayerNode struct {
	NodeLinks
}

func (node *LayerNode) Typeid() uint {
	return NODE_TYPE_RESIZE
}
func (node *LayerNode) GetArgs() string {
	args := fmt.Sprintf("-alpha set -channel a -evaluate set %d%", node.alpha)
	return args
}
