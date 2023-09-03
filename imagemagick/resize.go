package imagemagick

import "fmt"

type ResizeNode struct {
	filename string
	width    uint
	height   uint
	NodeLinks
}

func (node *ResizeNode) Typeid() uint {
	return NODE_TYPE_RESIZE
}
func (node *ResizeNode) GetArgs() string {
	args := fmt.Sprintf("-adaptive-resize %dx%d", node.width, node.height)
	return args
}

func (node *ResizeNode) SetImageResize(width, height uint) *ResizeNode {
	node.width = width
	node.height = height
	return node
}
