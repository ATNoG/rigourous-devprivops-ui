mod config;
mod dfd;
mod utils;

use std::panic;

use config::{FitConfig, Options};
use dfd::Data;

use wasm_bindgen::prelude::*;

#[wasm_bindgen]
extern "C" {
    fn alert(s: &str);

    // type Window;

    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);

    #[wasm_bindgen(js_name = "vis.Network")]
    type Network;

    #[wasm_bindgen(constructor)]
    fn new(container: &JsValue, data: &JsValue, options: &JsValue) -> Network;
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

fn log_jsval(v: &JsValue) {
    let serde_val = serde_wasm_bindgen::from_value::<serde_json::Value>(v.clone()).unwrap();
    log(serde_json::to_string(&serde_val).unwrap().as_str());
}

#[wasm_bindgen(start)]
pub fn module_initialization() {
    panic::set_hook(Box::new(console_error_panic_hook::hook));
}

#[wasm_bindgen]
pub fn process_and_render(s: &str) -> Result<(), JsValue> {
    let editor_content = Data::from_yaml(s);

    /*
    let data = serde_json::to_string(&editor_content).unwrap();

    let parsed_data: serde_json::Value = serde_json::from_str(&data)
        .map_err(|e| JsValue::from_str(&format!("Error parsing JSON: {}", e)))?;

    let nodes = parsed_data["nodes"].as_array().unwrap();
    let edges = parsed_data["edges"].as_array().unwrap();

    let document = web_sys::window().unwrap().document().unwrap();
    let container = document.get_element_by_id("mynetwork").unwrap();

    let nodes_data = js_sys::Array::new();
    for node in nodes {
        nodes_data.push(&serde_wasm_bindgen::to_value(&node).unwrap());
    }

    let edges_data = js_sys::Array::new();
    for edge in edges {
        edges_data.push(&serde_wasm_bindgen::to_value(&edge).unwrap());
    }
    */
    let nodes_data = js_sys::Array::new();
    for node in editor_content.nodes {
        nodes_data.push(&serde_wasm_bindgen::to_value(&node).unwrap());
    }

    let edges_data = js_sys::Array::new();
    for edge in editor_content.edges {
        edges_data.push(&serde_wasm_bindgen::to_value(&edge).unwrap());
    }

    let data = js_sys::Object::new();
    js_sys::Reflect::set(&data, &JsValue::from_str("nodes"), &nodes_data)?;
    js_sys::Reflect::set(&data, &JsValue::from_str("edges"), &edges_data)?;

    let opts = Options::new(true, true);
    let options = serde_wasm_bindgen::to_value(&opts)?;

    let document = web_sys::window().unwrap().document().unwrap();
    let container = document.get_element_by_id("mynetwork").unwrap();
    let network_constructor =
        js_sys::Reflect::get(&js_sys::global(), &JsValue::from_str("newNetwork"))?
            .dyn_into::<js_sys::Function>()?;
    let network =
        network_constructor.call3(&js_sys::global(), &container.into(), &data, &options)?;

    // Create a closure for the 'stabilizationIterationsDone' event
    let net = network.clone();
    let closure: Closure<dyn FnMut()> = Closure::new(move || {
        let fit_opts = FitConfig::new(500, String::from("easeInOutQuad"));

        let fit = js_sys::Reflect::get(&net, &JsValue::from_str("fit"))
            .unwrap()
            .dyn_into::<js_sys::Function>()
            .unwrap();
        let _ = fit
            .call1(&net, &serde_wasm_bindgen::to_value(&fit_opts).unwrap())
            .unwrap();
    });

    let closure_ref = closure.as_ref();
    let network_once = js_sys::Reflect::get(&network, &JsValue::from_str("once"))?
        .dyn_into::<js_sys::Function>()?;
    let _ = network_once.call2(
        &network,
        &JsValue::from_str("stabilizationIterationsDone"),
        closure_ref,
    )?;

    js_sys::Reflect::set(&js_sys::global(), &JsValue::from_str("network"), &network)?;

    closure.forget();
    Ok(())
}
