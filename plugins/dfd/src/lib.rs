mod dfd;
mod utils;

use dfd::{Data, Edge, Node};

use wasm_bindgen::prelude::*;

#[wasm_bindgen]
extern "C" {
    fn alert(s: &str);
}

#[wasm_bindgen]
pub fn greet(s: &str) -> String {
    s.into()
}

/*
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
*/

/*
#[wasm_bindgen(start)]
pub fn run_at_start() {
    alert("Hello World!");
}
*/

#[wasm_bindgen]
pub fn process_and_render(s: &str) -> Result<(), JsValue> {
    let editor_content = Data::from_yaml(s);

    let data = serde_json::to_string(&editor_content).unwrap();
    
    let parsed_data: serde_json::Value = serde_json::from_str(&data)
        .map_err(|e| JsValue::from_str(&format!("Error parsing JSON: {}", e)))?;

    let nodes = parsed_data["nodes"].as_array().unwrap();
    let edges = parsed_data["edges"].as_array().unwrap();

    let document = web_sys::window().unwrap().document().unwrap();
    let container = document.get_element_by_id("mynetwork").unwrap();

    // Assuming vis.js is available globally
    let nodes_data = js_sys::Array::new();
    for node in nodes {
        nodes_data.push(&JsValue::from_str(&node.to_string()));
    }

    let edges_data = js_sys::Array::new();
    for edge in edges {
        edges_data.push(&JsValue::from_str(&edge.to_string()));
    }

    let nodes_obj = js_sys::Object::new();
    js_sys::Reflect::set(&nodes_obj, &JsValue::from_str("nodes"), &nodes_data)?;

    let edges_obj = js_sys::Object::new();
    js_sys::Reflect::set(&edges_obj, &JsValue::from_str("edges"), &edges_data)?;

    let data = js_sys::Object::new();
    js_sys::Reflect::set(&data, &JsValue::from_str("nodes"), &nodes_obj)?;
    js_sys::Reflect::set(&data, &JsValue::from_str("edges"), &edges_obj)?;

    let options = js_sys::Object::new();
    js_sys::Reflect::set(
        &options,
        &JsValue::from_str("layout"),
        &JsValue::from_str("{ improvedLayout: true }"),
    )?;
    js_sys::Reflect::set(
        &options,
        &JsValue::from_str("physics"),
        &JsValue::from_str("{ enabled: true }"),
    )?;

    // Create the network (assuming vis.Network is available globally)
    let vis = js_sys::Reflect::get(&js_sys::global(), &JsValue::from_str("vis"))?
        .dyn_into::<js_sys::Object>()?;
    let network_constructor = js_sys::Reflect::get(&vis, &JsValue::from_str("Network"))?
        .dyn_into::<js_sys::Function>()?;

    network_constructor.call1(
        &vis,
        &js_sys::Array::of3(&container.into(), &data, &options),
    )?;

    // Optionally return the network instance or perform any further setup
    Ok(())
}
