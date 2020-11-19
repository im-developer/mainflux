CREATE TABLE users
(
	id bigserial constraint users_pk primary key,
	name varchar(255),
	secret varchar(255) not null,
	username varchar(255) not null,
	email varchar(255),
	password varchar(255),
	is_admin bool default false
);

CREATE UNIQUE index users_secret_uindex
	on users (secret);

CREATE UNIQUE index users_email_uindex
	on users (email);

CREATE UNIQUE index users_username_uindex
	on users (username);


