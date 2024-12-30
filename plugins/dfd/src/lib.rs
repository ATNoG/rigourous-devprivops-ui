mod config;
mod dfd;
mod utils;

use std::panic;

use config::Options;
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
        nodes_data.push(&serde_wasm_bindgen::to_value(&edge).unwrap());
    }

    let test_func = js_sys::Reflect::get(&js_sys::global(), &JsValue::from_str("test_func"))?
        .dyn_into::<js_sys::Function>()?;
    test_func.call1(&js_sys::global(), &JsValue::from_str("I work!"))?;

    let vis = js_sys::Reflect::get(&js_sys::global(), &JsValue::from_str("vis"))?
        .dyn_into::<js_sys::Object>()?;
    let window = web_sys::window().expect("no global `window` exists");

    /*
    let dataset_constructor =
        js_sys::Reflect::get(&vis, &JsValue::from_str("DataSet"))?.dyn_into::<js_sys::Object>()?;

    let nodes_ds = dataset_constructor.constructor().call1(&vis, &nodes_data)?;
    let edges_ds = dataset_constructor.constructor().call1(&vis, &edges_data)?;

    log_jsval(&nodes_ds);
    log_jsval(&edges_ds);
    */

    /*
    let nodes_obj = js_sys::Object::new();
    js_sys::Reflect::set(&nodes_obj, &JsValue::from_str("nodes"), &nodes_data)?;

    let edges_obj = js_sys::Object::new();
    js_sys::Reflect::set(&edges_obj, &JsValue::from_str("edges"), &edges_data)?;
    */

    let data = js_sys::Object::new();
    js_sys::Reflect::set(&data, &JsValue::from_str("nodes"), &nodes_data)?;
    js_sys::Reflect::set(&data, &JsValue::from_str("edges"), &edges_data)?;

    /*
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
    */

    let opts = Options::new(true, true);
    let options = serde_wasm_bindgen::to_value(&opts)?;

    // Create the network (assuming vis.Network is available globally)
    // let network_constructor = js_sys::Reflect::get(&vis, &JsValue::from_str("Network"))?
    //     .dyn_into::<js_sys::Function>()?;
    // let network =
    //     js_sys::Reflect::get(&vis, &JsValue::from_str("Network"))?.dyn_into::<js_sys::Function>()?;
    // let _network_constructor = network.constructor();

    /*
    let res = network.call3(
        &network,
        &container.into(),
        &data,
        &options,
    )?; // <-----
    */
    // let res = Network::new(&container.into(), &data, &options);
    let network_constructor =
        js_sys::Reflect::get(&js_sys::global(), &JsValue::from_str("newNetwork"))?
            .dyn_into::<js_sys::Function>()?;
    let network =
        network_constructor.call3(&js_sys::global(), &container.into(), &data, &options)?;

    alert("here 5!");
    // let _ = Network::new(&container.into(), &data, &options);

    log_jsval(&network);

    // Create a closure for the 'stabilizationIterationsDone' event
    let closure = Closure::wrap(Box::new(move || {
        // Fit the graph with animation options once stabilized
        let fit_opts = js_sys::Object::new();
        let animation_opts = js_sys::Object::new();
        js_sys::Reflect::set(&animation_opts, &"duration".into(), &500.into()).unwrap();
        js_sys::Reflect::set(
            &animation_opts,
            &"easingFunction".into(),
            &"easeInOutQuad".into(),
        )
        .unwrap();
        js_sys::Reflect::set(&fit_opts, &"animation".into(), &animation_opts).unwrap();

        // Assuming `network` is a global variable for `vis.Network`
        // let network: js_sys::Object = window.get("network").unwrap().dyn_into().unwrap();
        let fit = js_sys::Reflect::get(&network, &JsValue::from_str("fit"))
            .unwrap()
            .dyn_into::<js_sys::Function>()
            .unwrap();
        let _: JsValue = fit.call1(&network, &fit_opts).unwrap();
    }) as Box<dyn FnMut()>);

    // Assign closure to the network event
    let closure_ref = closure.as_ref();
    // Add the event listener (using the vis.js event listener interface)
    let network = js_sys::Object::new(); // Here you'd instantiate your vis.Network
    /*
    let _: JsValue = js_sys::Reflect::apply(
        &network,
        &"once".into(),
        &js_sys::Array::of2(&"stabilizationIterationsDone".into(), &closure_ref),
    )
    .unwrap();
    */

    // Set the network as a global variable (like in the JavaScript code)
    // window.set("network", network);

    // Optionally return the network instance or perform any further setup
    Ok(())
}

/*
To translate this JavaScript code into Rust code using wasm-bindgen for WebAssembly, we need to interact with the JavaScript vis.Network API in Rust. This requires setting up the wasm-bindgen crate to bridge between Rust and JavaScript, and using JavaScript's vis.js library.

Here is a step-by-step translation:

    Cargo Setup: In your Cargo.toml, include the wasm-bindgen dependency to create bindings between Rust and JavaScript.

    [dependencies]
    wasm-bindgen = "0.2"
    js-sys = "0.3"
    web-sys = "0.3"

    Rust Code with wasm-bindgen: You'll need to use wasm-bindgen to access JavaScript's vis.Network API from Rust. You can use web_sys to interact with the DOM (for container) and JavaScript objects.

    Here’s how you can translate your JavaScript code to Rust using wasm-bindgen:

```rs
use wasm_bindgen::prelude::*;
use web_sys::{Document, Window, Element};
use js_sys::Function;
use wasm_bindgen::JsCast;

#[wasm_bindgen(start)]
pub fn start() -> Result<(), JsValue> {
    // Initialize the wasm-bindgen logger (optional, for debugging purposes)
    console_error_panic_hook::set_once();

    // Access the window and document
    let window: Window = web_sys::window().expect("no global `window` exists");
    let document: Document = window.document().expect("should have a document on window");

    // Get the container element where the graph will be rendered
    let container = document.get_element_by_id("container")
        .expect("container element not found");

    // Create the 'data' and 'options' for the Network (These can be custom Rust types or directly set up in JS)
    // Here, we assume `data` and `options` are JavaScript objects, and we're passing them from Rust to JS.

    // Setup a simple network initialization
    let vis = js_sys::Object::new();  // This represents the vis.js Network object

    // Create a closure for the 'stabilizationIterationsDone' event
    let closure = Closure::wrap(Box::new(move || {
        // Fit the graph with animation options once stabilized
        let fit_opts = js_sys::Object::new();
        let animation_opts = js_sys::Object::new();
        js_sys::Reflect::set(&animation_opts, &"duration".into(), &500.into()).unwrap();
        js_sys::Reflect::set(&animation_opts, &"easingFunction".into(), &"easeInOutQuad".into()).unwrap();
        js_sys::Reflect::set(&fit_opts, &"animation".into(), &animation_opts).unwrap();

        // Assuming `network` is a global variable for `vis.Network`
        let network: js_sys::Object = window.get("network").unwrap().dyn_into().unwrap();
        let _: JsValue = js_sys::Reflect::apply(&network, &JsValue::NULL, &js_sys::Array::of1(&fit_opts)).unwrap();
    }) as Box<dyn FnMut()>);

    // Assign closure to the network event
    let closure_ref = closure.as_ref();
    // Add the event listener (using the vis.js event listener interface)
    let network = js_sys::Object::new(); // Here you'd instantiate your vis.Network
    let _: JsValue = js_sys::Reflect::apply(&network, &"once".into(), &js_sys::Array::of2(&"stabilizationIterationsDone".into(), &closure_ref)).unwrap();

    // Set the network as a global variable (like in the JavaScript code)
    window.set("network", network);

    Ok(())
}
```

Key Points:

    wasm-bindgen: We use wasm-bindgen to bind Rust with JavaScript functions (like vis.Network).
    Closure: Rust closures are used for handling the event listeners, such as stabilizationIterationsDone.
    Accessing the DOM: We access the DOM (where the container element is) via web_sys::document, similar to how you'd do it in JavaScript.
    Global Variable: We set a network global variable to access it outside the Rust environment.

Missing Components:

    data and options: These should be created or passed from Rust to JavaScript as JavaScript objects. You can create them manually in JavaScript or use the Rust JsValue and js_sys types to create such objects.
    container element: The container in JavaScript is assumed to be an HTML element that must exist in the DOM.

Compiling and Running:

To compile this Rust code into WebAssembly, you would typically use wasm-pack or cargo build --target wasm32-unknown-unknown to build the WebAssembly module. You would also need an HTML/JS frontend that loads this Wasm module and binds it to the appropriate DOM elements.

The JavaScript part (outside of Rust) would look like this:

```js
import * as wasm from './pkg/your_project_name';

window.onload = () => {
  wasm.start(); // Initializes the Rust WebAssembly code
};
```

This provides the structure to interact with JavaScript's vis.js library from Rust using WebAssembly!

*/
