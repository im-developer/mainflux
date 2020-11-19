pub mod model;
pub mod schema;
pub mod auth;

use rocket_contrib::json::{Json, JsonValue};
use self::auth::crypto::sha2::Sha256;
use self::auth::jwt::{
    Header,
    Registered,
    Token,
};
use crate::db;
use crate::user::model::User;
use rocket::http::Status;
use crate::user::auth::ApiKey;

#[derive(Serialize, Deserialize)]
struct Credentials {
    username: String,
    password: String
}

// #[get("/<id>")]
// fn read_one(id: i32, connection: db::Connection) -> Result<Json<JsonValue>, Status> {
//     User::read(id, &connection)
//         .map(|item| Json(json!(item)))
//         .map_err(|_| Status::NotFound)
// }

#[get("/json")]
fn read_oxne() -> JsonValue {
    json!({
        "id": 83,
        "values": [1, 2, 3, 4]
    })
}

#[get("/sensitive")]
fn sensitive(key: ApiKey) -> String {
    format!("Hello, you have been identified as {}", key.0)
}


#[post("/login", format = "json", data = "<credentials>")]
fn login(credentials: Json<Credentials>, connection: db::Connection) -> Result<Json<JsonValue>, Status> {
    let header: Header = Default::default();
    let username = credentials.username.to_string();
    let password = credentials.password.to_string();

    match User::by_username_and_password(username, password, &connection) {
        None => {
            Err(Status::NotFound)
        },
        Some(user) => {
            let claims = Registered {
                sub: Some(user.name.into()),
                ..Default::default()
            };
            let token = Token::new(header, claims);

            token.signed(b"secret_key", Sha256::new())
                .map(|message| Json(json!({ "success": true, "token": message })))
                .map_err(|_| Status::InternalServerError)
        }
    }
}

pub fn mount(rocket: rocket::Rocket) -> rocket::Rocket {
    rocket
    //     .mount("/user", routes![read, read_error, read_one, create, update, delete, info, info_error])
        .mount("/auth", routes![login, sensitive])
}
