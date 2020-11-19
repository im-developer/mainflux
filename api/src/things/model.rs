use diesel;
use diesel::prelude::*;
use diesel::pg::PgConnection;
use crate::things::schema::things;
use crate::things::schema::channels;
use diesel::sql_types::Uuid;

#[derive(Serialize, Identifiable, Clone, Queryable, Debug)]
pub struct Thing {
    pub id: Uuid
}

#[derive(Serialize, Identifiable, Clone, Queryable, Debug)]
pub struct Channel {
    pub id: String,
    pub owner: String,
    pub name: String,
    pub metadata: serde_json::Value
}

impl Thing {

    pub fn read(id: Uuid, connection: &PgConnection) -> QueryResult<Vec<Thing>> {
        // println!("{}", id);
        //if id != 0 {
            things::table.find(id).load::<Thing>(connection)
        //} else {
        //    things::table.order(things::id).load::<Thing>(connection)
        //}
    }

    pub fn get_one(id: Uuid, connection: &PgConnection) -> QueryResult<Thing> {

        let results = things::table
            .limit(5)
            .load::<Thing>(connection)
            .expect("Error loading posts");

        // println!("Displaying {} posts", results.len());
        // for post in results {
        //     println!("{}", post.id);
        //     println!("----------\n");
        // }

        // println!("{}", id);
        things::table.find(id).first(connection)
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
