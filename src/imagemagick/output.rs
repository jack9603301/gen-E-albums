use super::global::NODE_TYPE_OUTPUT;
use super::node::Node;
use derives::node_manager;

#[node_manager]
pub struct OutputNode {
    filename : String
}

impl Node for OutputNode {
    fn type_id(& self) -> u32 {
        return NODE_TYPE_OUTPUT;
    }
    fn get_args(self: & OutputNode) -> Option<String> {
        return match self.filename.len() {
            0 => None,
            _ => Some(self.filename.clone())
        };
    }
}