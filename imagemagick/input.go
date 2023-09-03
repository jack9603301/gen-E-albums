package imagemagick

import "fmt"

const INPUT_NODE_OPENED uint = 1
const INPUT_NODE_CREATE uint = 2

type InputNodeOpenedParam struct {
	filename string
}

type InputNodeCreateParam struct {
	height uint
	width  uint
}

type InputNode struct {
	arg       interface{}
	inputType uint
	NodeLinks
}

func (node *InputNode) Typeid() uint {
	return NODE_TYPE_INPUT
}

func (node *InputNode) GetArgs() string {
	var args string
	switch node.inputType {
	case INPUT_NODE_OPENED:
		open_param := node.arg.(InputNodeOpenedParam)
		args = fmt.Sprintf("%s", open_param.filename)
	case INPUT_NODE_CREATE:
		create_param := node.arg.(InputNodeCreateParam)
		args = fmt.Sprintf("xc:none -resize %dx%d", create_param.width, create_param.height)
	}

	return args
}

func SetFileName(node *Node, filename string) *Node {
	type_id := (*node).Typeid()
	switch type_id {
	case NODE_TYPE_INPUT:
		input_node := (*node).(*InputNode)
		arg := &InputNodeOpenedParam{}
		arg.filename = filename
		input_node.arg = arg
		input_node.inputType = INPUT_NODE_OPENED
	case NODE_TYPE_OUTPUT:
		output_node := (*node).(*OutputNode)
		output_node.filename = filename
	}

	return node
}

func CreateImage(node *Node, width, height uint) *Node {
	type_id := (*node).Typeid()
	switch type_id {
	case NODE_TYPE_INPUT:
		input_node := (*node).(InputNode)
		arg := &InputNodeCreateParam{}
		arg.width = width
		arg.height = height
		input_node.arg = arg
		input_node.inputType = INPUT_NODE_CREATE
	default:
		return nil
	}

	return node
}
