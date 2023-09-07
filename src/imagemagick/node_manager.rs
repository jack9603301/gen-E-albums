

use std::rc::Rc;
use std::str::FromStr;

use super::node::Node;
use super::input::InputNode;

pub struct NodeManager {
    process_list: Vec<Rc<dyn Node>>,
    current_item: i32
}

impl NodeManager {
    pub fn new() -> NodeManager {
        return NodeManager{
            process_list: Vec::new(),
            current_item: -1
        };
    }
    pub fn get_len(&self) -> usize {
        assert!((self.current_item + 1) as usize == self.process_list.len() as usize);
        return self.process_list.len();
    }
    pub fn input_open(&mut self, filename: String) -> &mut NodeManager {
        assert!(self.get_len() == 0);
        let mut input_node = InputNode::new();
        input_node.set_file_name(filename);
        let mut_node = Rc::new(input_node);
        self.process_list.push(mut_node);
        self.current_item += 1;
        return self;
    }
    pub fn input_create(&mut self, width: i32, height: i32) -> &mut NodeManager {
        assert!(self.get_len() == 0);
        let mut input_node = InputNode::new();
        input_node.set_create_empty_image(width, height);
        let mut_node = Rc::new(input_node);
        self.process_list.push(mut_node);
        self.current_item += 1;
        return self;
    }
    pub fn run(&self) -> String {
        let mut result = String::from("magick convert ");
        let list = &self.process_list;
        for node in list {
            result += &node.get_args().unwrap();
            result += &String::from(" ");
        }
        return result;
    }
}