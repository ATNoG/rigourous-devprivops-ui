mod utils;
mod dfd;

use dfd::{Node, Edge, Data};

use wasm_bindgen::prelude::*;

#[wasm_bindgen]
extern "C" {
    fn alert(s: &str);
}

#[wasm_bindgen]
pub fn greet(s: &str) -> String {
    s.into()
}

#[wasm_bindgen]
pub fn process(s: &str) -> String {
    /*
    let n1 = Node::new(String::from("1"), s.into());
    let n2 = Node::new(String::from("2"), "a".into());

    let e1 = Edge::new(n1.get_id(), n2.get_id());

    let data = Data::new(vec![n1, n2], vec![e1]);
    let json_res = serde_json::to_string(&data).unwrap();

    json_res
    */
    let data = Data::from_yaml(s);

    let json_res = serde_json::to_string(&data).unwrap();

    json_res
}
