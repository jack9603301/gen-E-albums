extern crate proc_macro;
use proc_macro::TokenStream;
use quote::quote;
use syn::parse::Parser;
use syn::{parse_macro_input, DeriveInput};

#[proc_macro_attribute]
pub fn node_manager(_attr: proc_macro::TokenStream, input: proc_macro::TokenStream) -> TokenStream {

    let mut ast = parse_macro_input!(input as DeriveInput);   //解析输入token序列，转化成语法树

    let name : & syn::Ident = &ast.ident;

    let node_next= quote! {
        NextNode: Option<Rc<dyn Node>>
    };  //定义NextNode的AST（语法树）描述
    let node_prev= quote! {
        PrevNode: Option<Rc<dyn Node>>
    };//定义PrevNode的AST（语法树）描述

    let node_next_ast = syn::Field::parse_named.parse2(node_next).unwrap();   //解析，生成node_next的代码对应的AST
    let node_prev_ast = syn::Field::parse_named.parse2(node_prev).unwrap();   //解析，生成node_prev的代码对应的AST

    match &mut ast.data {
        syn::Data::Struct(ref mut struct_data) => {
            match &mut struct_data.fields {
                syn::Fields::Named(fields) => {
                    fields.named.push(node_next_ast);
                    fields.named.push(node_prev_ast);
                }
                _ => {
                    ()
                }
            }
        }
        _ => {
            ()
        }
    }   //将这两个属性添加到AST中，代码转换

    //生成输出的标记

    let output = quote! {
        #ast

        impl NodeManager for #name {
            fn GetNextNode(&self) -> &Option<Rc<dyn Node>> {
                return &self.NextNode;
            }
            fn GetPrevNode(&self) -> &Option<Rc<dyn Node>> {
                return &self.PrevNode;
            }
            fn SetNextNode(&mut self, node : &Option<Rc<dyn Node>>) {
                self.NextNode = match node {
                    Some(node) => Some(node.clone()),
                    None => None
                };
            }
            fn SetPrevNode(&mut self, node : &Option<Rc<dyn Node>>) {
                self.PrevNode = match node {
                    Some(node) => Some(node.clone()),
                    None => None
                }
            }
        }
    };

    TokenStream::from(output)
    
}