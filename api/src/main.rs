#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate rocket;
#[macro_use] extern crate diesel;
#[macro_use] extern crate serde_derive;
#[macro_use] extern crate rocket_contrib;

mod db;
// mod user;
mod store;
mod things;

use diesel::prelude::*;
// use user::model::User;

fn main() {
    // error[E0425]: cannot find value `static_rocket_route_info_for_world` in this scope
    // rocket::ignite().mount("/", routes![hello, param, other::world]).launch();
    let mut rocket = rocket::ignite()
        .manage(db::connect());
    // rocket = user::mount(rocket);
    rocket = things::mount(rocket);
    rocket = store::mount(rocket);
    rocket.launch();

}
