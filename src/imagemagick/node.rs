use std::rc::Rc;

pub trait Node {
    fn type_id(& self) -> u32;
    fn get_args(& self) -> Option<String>;
}

pub trait NodeManager {
    fn set_next_node(&mut self, node : &Option<Rc<dyn Node>>); 
    fn get_next_node(&self) -> &Option<Rc<dyn Node>>;
    fn set_prev_node(&mut self, node : &Option<Rc<dyn Node>>);
    fn get_prev_node(&self) -> &Option<Rc<dyn Node>>;
}