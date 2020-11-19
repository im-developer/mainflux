extern crate serde_json;
table! {
    stores (id) {
        id -> BigInt,
        package_id -> BigInt,
        currency_id -> BigInt,
        language_id -> BigInt,
        maintenance -> Bool,
        show_sold_out_products -> Bool,
        display_testimonials -> Bool,
        theme -> Varchar,
        url -> Varchar,
        logo -> Varchar,
        name -> Jsonb,
    }
}
