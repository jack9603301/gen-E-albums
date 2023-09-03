package imagemagick

type Image struct {
	NodeInfo map[uint][]*Node
}

func (image *Image) Input(filename string) *Node {
	var new_node Node
	new_node = &InputNode{}
	node := SetFileName(&new_node, filename)
	image.NodeInfo[NODE_TYPE_INPUT] = append(image.NodeInfo[NODE_TYPE_INPUT], node)
	return node
}
