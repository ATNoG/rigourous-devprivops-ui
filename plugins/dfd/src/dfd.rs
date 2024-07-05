use serde::{Serialize,Deserialize};

#[derive(Serialize,Deserialize)]
pub struct Node {
    id: usize,
    label: String,
}

#[derive(Serialize,Deserialize)]
pub struct Edge {
    from: usize,
    to: usize,
}

#[derive(Serialize,Deserialize)]
pub struct Data {
    nodes: Vec<Node>,
    edges: Vec<Edge>,
}

impl Node {
    pub fn new(id: usize, label: String) -> Self {
        Node {id, label}
    }

    pub fn get_id(&self) -> usize {
        self.id
    }
}

impl Edge {
    pub fn new(from: usize, to: usize) -> Self {
        Edge {from, to}
    }
}

impl Data {
    pub fn new(nodes: Vec<Node>, edges: Vec<Edge>) -> Self {
        Data {nodes, edges}
    }
}   

