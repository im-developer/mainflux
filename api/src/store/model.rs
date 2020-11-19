use diesel;
use diesel::prelude::*;
use diesel::pg::PgConnection;
use crate::store::schema::stores;

#[derive(Serialize, Identifiable, Clone, Queryable, Debug)]
pub struct Store {
    pub id: i64,
    pub package_id: i64,
    pub currency_id: i64,
    pub language_id: i64,
    pub maintenance: bool,
    pub show_sold_out_products: bool,
    pub display_testimonials: bool,
    pub theme: String,
    pub url: String,
    pub logo: String,
    pub name: serde_json::Value
}

impl Store {

    pub fn read(id: i64, connection: &PgConnection) -> QueryResult<Vec<Store>> {
        if id != 0 {
            stores::table.find(id).load::<Store>(connection)
        } else {
            stores::table.order(stores::id).load::<Store>(connection)
        }
    }

    pub fn get_one(id: i64, connection: &PgConnection) -> QueryResult<Store> {
        stores::table.find(id).first(connection)
    }

    // pub fn by_username_and_password(username_: String, password_: String, connection: &MysqlConnection) -> Option<User> {
    //     let res = users::table
    //         .filter(users::name.eq(username_))
    //         .filter(users::password.eq(password_))
    //         .order(users::id)
    //         .first(connection);
    //     match res {
    //         Ok(user) => Some(user),
    //         Err(_) => {
    //             None
    //         }
    //     }
    // }

    // pub fn update(id: i32, user: User, connection: &MysqlConnection) -> bool {
    //     diesel::update(users::table.find(id)).set(&user).execute(connection).is_ok()
    // }
    //
    // pub fn delete(id: i32, connection: &MysqlConnection) -> bool {
    //     diesel::delete(users::table.find(id)).execute(connection).is_ok()
    // }
}
