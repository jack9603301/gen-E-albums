use std::rc::Rc;

pub trait Node {
    fn TypeId(& self) -> u32;
    fn GetArgs(&mut self) -> Option<String>;
}

pub trait NodeManager {
    fn SetNextNode(&mut self, node : &Option<Rc<dyn Node>>); 
    fn GetNextNode(&self) -> &Option<Rc<dyn Node>>;
    fn SetPrevNode(&mut self, node : &Option<Rc<dyn Node>>);
    fn GetPrevNode(&self) -> &Option<Rc<dyn Node>>;
}