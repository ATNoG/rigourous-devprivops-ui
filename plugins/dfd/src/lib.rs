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
    let n1 = Node::new(1, s.into());
    let n2 = Node::new(2, "a".into());

    let e1 = Edge::new(n1.get_id(), n2.get_id());

    let data = Data::new(vec![n1, n2], vec![e1]);
    let json_res = serde_json::to_string(&data).unwrap();

    json_res
    /*
    r#" 
    {
        "nodes": [
            { "id": 1, "label": "Node 1" },
            { "id": 2, "label": "Node 2" },
            { "id": 3, "label": "Node 3" },
            { "id": 4, "label": "Node 4" },
            { "id": 5, "label": "Node 5" }
        ],
        "edges": [
            { "from": 1, "to": 2 },
            { "from": 1, "to": 3 },
            { "from": 2, "to": 4 },
            { "from": 2, "to": 5 }
        ]
    }
    "#.into()
    */
}
