extern crate serde_json;

table! {
    things (owner) {
        owner -> Varchar,
    }
}

table! {
    channels (id) {
        id -> Varchar,
        owner -> Varchar,
        name -> Varchar,
        metadata -> Jsonb,
    }
}
