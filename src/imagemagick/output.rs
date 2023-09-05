use super::global::NODE_TYPE_OUTPUT;
use super::node::Node;

pub struct OutputNode {
    filename : String
}

impl Node for OutputNode {
    fn TypeId(& self) -> u32 {
        return NODE_TYPE_OUTPUT;
    }
}