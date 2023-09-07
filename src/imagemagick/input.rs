extern crate derives;
use std::rc::Rc;
use std::string::String;

use derives::node_manager;
use super::global::NODE_TYPE_INPUT;
use super::node::Node;


#[derive(Clone)]
struct OpenInputType {
    filename: Option<Rc<String>>
}

#[derive(Clone)]
struct CreateInputType {
    width : i32,
    height : i32
}

#[derive(Copy, Clone)]
#[derive(PartialEq)]
pub enum InputTypeE {
    None,
    Open,
    Create
}

#[derive(Clone)]
struct InputParam {
    open : OpenInputType,
    create : CreateInputType
}

#[node_manager]
#[derive(Clone)]
pub struct InputNode {
    iparam : InputParam,
    input_type : InputTypeE
}

impl InputNode {
    pub fn new() -> InputNode {
        let node = InputNode{
            iparam: InputParam{
                create: CreateInputType{
                    width: -1,
                    height: -1
                },
                open: OpenInputType{
                    filename: None
                }
            },
            input_type: InputTypeE::None,
            next_node: None,
            prev_node: None,
        };
        return node
    }
    pub fn get_input_type(& self) -> InputTypeE {
        return self.input_type;
    }
    pub fn set_input_type(&mut self, input_type : InputTypeE) {
        self.input_type = input_type;
    }
    pub fn set_file_name(&mut self, filename: String) {
        let filename_str: String = String::from(filename);
        let str = Rc::new(filename_str);
        self.iparam.open.filename = Some(str);
        self.iparam.create.width = -1;
        self.iparam.create.height = -1;
        self.set_input_type(InputTypeE::Open);
    }
    pub fn set_create_empty_image(&mut self, width: i32, height: i32) {
        self.iparam.create.width = width;
        self.iparam.create.height = height;
        self.iparam.open.filename = None;
        self.set_input_type(InputTypeE::Create);
    }
}

impl Node for InputNode {
    fn type_id(& self) -> u32 {
        return NODE_TYPE_INPUT;
    }
    fn get_args(& self) -> Option<String> {
        if self.get_input_type() == InputTypeE::Open {
            let filename = match &self.iparam.open.filename {
                Some(filename) => filename.to_string(),
                None => String::from("")
            };
            return Some(format!("-i {}", filename));
        } else if self.get_input_type() == InputTypeE::Create {
            return Some(format!("-i xc:none -resize {}x{}", self.iparam.create.width, self.iparam.create.height));
        } else {
            return None;
        }
    }
}