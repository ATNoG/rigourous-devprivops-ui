use serde::{Serialize,Deserialize};
use serde_yml::Value;
use core::panic;
use std::collections::HashMap;

#[derive(Serialize,Deserialize)]
pub struct Node {
    id: String,
    label: String,
}

#[derive(Serialize,Deserialize)]
pub struct Edge {
    from: String,
    to: String,
}

#[derive(Serialize,Deserialize)]
pub struct Data {
    nodes: Vec<Node>,
    edges: Vec<Edge>,
}

impl Node {
    pub fn new(id: String, label: String) -> Self {
        let mut id = id;
        if id.starts_with(":") {
            id = id[1..].to_string()
        }

        Node { id, label, }
    }
}

impl Edge {
    pub fn new(from: String, to: String) -> Self {
        let mut from = from;
        if from.starts_with(":") {
            from = from[1..].to_string()
        }
        
        let mut to = to;
        if to.starts_with(":") {
            to = to[1..].to_string()
        }

        Edge {from, to}
    }
}

impl Data {
    pub fn new(nodes: Vec<Node>, edges: Vec<Edge>) -> Self {
        Data {nodes, edges}
    }

    pub fn from_yaml(yaml_str: &str) -> Data {
        let mut nodes = vec![];
        let mut edges = vec![];

        let yaml: HashMap<String, Value> = serde_yml::from_str(yaml_str).expect("Could not deserialize yaml");

        let data_flows = &yaml["data flows"];
        if let Value::Sequence(df) = data_flows {
            for v in df {
                if let Value::Mapping(m) = v {

                    /*
                    let id = if let Some(Value::String(value)) = m.get(Value::String(String::from("id"))) {
                        value.clone()
                    } else { panic!("No id") };
                     */
                    let to = if let Some(Value::String(value)) = m.get(Value::String(String::from("to"))) {
                        value.clone()
                    } else { panic!("No to") };
                    let from = if let Some(Value::String(value)) = m.get(Value::String(String::from("from"))) {
                        value.clone()
                    } else { panic!("No from") };

                    // nodes.push(Node::new(id.clone(), id));
                    edges.push(Edge::new(from, to))
                }
            }
        } else { panic!("Could not parse data flows") }

        let external_entities = &yaml["external entities"];
        if let Value::Sequence(ee) = external_entities {
            for v in ee {
                if let Value::Mapping(m) = v {

                    let id = if let Some(Value::String(value)) = m.get(Value::String(String::from("id"))) {
                        value.clone()
                    } else { panic!("No id") };

                    nodes.push(Node::new(id.clone(), id));
                }
            }
        } else { panic!("Could not parse external entities") }

        let processes = &yaml["processes"];
        if let Value::Sequence(pr) = processes {
            for v in pr {
                if let Value::Mapping(m) = v {

                    let id = if let Some(Value::String(value)) = m.get(Value::String(String::from("id"))) {
                        value.clone()
                    } else { panic!("No id") };

                    nodes.push(Node::new(id.clone(), id));
                }
            }
        } else { panic!("Could not parse processes") }

        let data_stores = &yaml["data stores"];
        if let Value::Sequence(ds) = data_stores {
            for v in ds {
                if let Value::Mapping(m) = v {

                    let id = if let Some(Value::String(value)) = m.get(Value::String(String::from("id"))) {
                        value.clone()
                    } else { panic!("No id") };

                    nodes.push(Node::new(id.clone(), id));
                }
            }
        } else { panic!("Could not parse data stores") }

        Data { nodes, edges }
    }
}   

