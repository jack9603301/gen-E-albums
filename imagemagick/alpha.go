package imagemagick

import "fmt"

type AlphaNode struct {
	alpha uint
	NodeLinks
}

func (node *AlphaNode) Typeid() uint {
	return NODE_TYPE_RESIZE
}
func (node *AlphaNode) GetArgs() string {
	args := fmt.Sprintf("-alpha set -channel a -evaluate set %d%", node.alpha)
	return args
}

func (node *AlphaNode) SetImageAlpha(alpha uint) *AlphaNode {
	node.alpha = alpha
	return node
}
