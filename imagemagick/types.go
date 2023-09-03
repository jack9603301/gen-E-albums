package imagemagick

const NODE_TYPE_INPUT uint = 1
const NODE_TYPE_OUTPUT uint = 2
const NODE_TYPE_RESIZE uint = 3
const NODE_TYPE_LAYER uint = 4
const NODE_TYPE_ALPHA uint = 5

type KwArgs_Type map[string]interface{}

type Node interface {
	Typeid() uint
	GetArgs() string
}

type NodeLinks struct {
	PrevNode []*Node
	NextNode *Node
}

// 将指定节点连接到此节点的下级
func (link *NodeLinks) LinkNode(node *Node) bool {
	type_id := (*node).Typeid()
	switch type_id {
	case NODE_TYPE_INPUT:
		intput_next_node := (*node).(*InputNode)
		link.NextNode = node
		intput_next_node.PrevNode = append(intput_next_node.PrevNode, node)
	case NODE_TYPE_OUTPUT:
		output_next_node := (*node).(*InputNode)
		link.NextNode = node
		output_next_node.PrevNode = append(output_next_node.PrevNode, node)
	default:
		return false
	}
	return false
}
