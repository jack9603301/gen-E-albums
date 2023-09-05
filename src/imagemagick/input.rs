extern crate derives;
use std::rc::Rc;
use std::string::String;

use derives::node_manager;
use super::global::NODE_TYPE_INPUT;
use super::node::Node;
use super::node::NodeManager;

#[derive(Clone)]
struct OpenInputType {
    filename: Rc<String>
}

struct CreateInputType {
    width : i32,
    height : i32
}

#[derive(Copy, Clone)]
#[derive(PartialEq)]
pub enum InputTypeE {
    Open,
    Create
}

struct InputParam {
    open : OpenInputType,
    create : CreateInputType
}

#[node_manager]
pub struct InputNode {
    iparam : InputParam,
    InputType : InputTypeE
}

impl InputNode {
    pub fn GetInputType(&mut self) -> InputTypeE {
        return self.InputType;
    }
    pub fn SetInputType(&mut self, input_type : InputTypeE) {
        self.InputType = input_type;
    }
}

impl Node for InputNode {
    fn TypeId(& self) -> u32 {
        return NODE_TYPE_INPUT;
    }
    fn GetArgs(&mut self) -> Option<String> {
        if self.GetInputType() == InputTypeE::Open {
            let filename = self.iparam.open.filename.to_string();
            return Some(format!("-i {}", filename));
        } else if self.GetInputType() == InputTypeE::Create {
            return Some(format!("-i xc:none -resize {}x{}", self.iparam.create.width, self.iparam.create.height));
        } else {
            return None;
        }
    }
}