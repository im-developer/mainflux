pub mod model;
pub mod schema;

use rocket_contrib::json::{Json, JsonValue};
use crate::db;
use rocket::http::Status;
use crate::things::model::Thing;
use diesel::sql_types::Uuid;

#[get("/<id>")]
fn read_one(id: uuid::Uuid, connection: db::Connection) -> Result<Json<JsonValue>, Status> {
    Thing::get_one(id, &connection)
        .map(|item| Json(json!(item)))
        .map_err(|_| Status::NotFound)
}

pub fn mount(rocket: rocket::Rocket) -> rocket::Rocket {
    rocket
        .mount("/things", routes![read_one])
}
