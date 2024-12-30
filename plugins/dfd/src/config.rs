use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct Options {
    layout: LayoutOpts,
    physics: PhysicsOpts,
}

impl Options {
    pub fn new(improved_layout: bool, enabled: bool) -> Self {
        return Self {
            layout: LayoutOpts { improved_layout },
            physics: PhysicsOpts { enabled },
        };
    }
}

#[derive(Serialize, Deserialize)]
struct LayoutOpts {
    #[serde(rename = "improvedLayout")]
    improved_layout: bool,
}

#[derive(Serialize, Deserialize)]
struct PhysicsOpts {
    enabled: bool,
}
