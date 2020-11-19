pub mod model;
pub mod schema;

use rocket_contrib::json::{Json, JsonValue};
use crate::db;
use rocket::http::Status;
use crate::store::model::Store;

#[get("/<id>")]
fn read_one(id: i64, connection: db::Connection) -> Result<Json<JsonValue>, Status> {
    Store::get_one(id, &connection)
        .map(|item| Json(json!(item)))
        .map_err(|_| Status::NotFound)
}

pub fn mount(rocket: rocket::Rocket) -> rocket::Rocket {
    rocket
        .mount("/stores", routes![read_one])
}
